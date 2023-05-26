// Author: huaxr
// Time:   2022/1/4 上午10:29
// Git:    huaxr

package core

import (
	"container/heap"
	ctx2 "context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/spf13/cast"

	"github.com/huaxr/magicflow/component/dispatch"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/pkg/confutil"
)

var (
	dispatcher *_dispatcher
	semaphore  = make(chan struct{}, 1e6)
)

type _dispatcher struct {
	lock sync.Mutex

	ctx ctx2.Context
	// for the sake of refrain from trigger message rapid backlog,
	// using the msgQ which implement container/heap in order
	// to providing a priority supportive scheme or sync message handle.
	msgQ heapQ
	// event chan
	schemeC chan scheme
	// record db
	recorder *recorder
	// handler ToServer&ToTrigger will convert internal message and requeue to the broker.
	handler func(msg *message) (err error)

	// publisher
	publisher dispatch.Publisher
}

func LaunchCore(ctx ctx2.Context, mq dispatch.MQ) {
	dispatcher = new(_dispatcher)
	dispatcher.lock = sync.Mutex{}
	dispatcher.ctx = ctx
	dispatcher.msgQ = &priorityMq{}
	heap.Init(dispatcher.msgQ)
	dispatcher.schemeC = make(chan scheme, 1e6)
	dispatcher.recorder = initDaoJob(ctx)
	dispatcher.publisher = dispatch.InitMQ(ctx, mq)

	dispatcher.handler = func(msg *message) (err error) {
		if gCache.HasPlaybook(msg.Task.getPlaybookId()) {
			hatches := msg.Hatch(true)
			if err = dispatcher.dispatchMany(hatches); err != nil {
				return
			}
			return
		} else {
			err = fmt.Errorf("playbook not found %v", msg.Task.getPlaybookId())
			return
		}
	}

	launchExchangeCache()
	launchG()

	thread := confutil.GetConf().GetDispatchThreadCount()
	logx.L().Infof("dispatcher threads count: %d", thread)
	for i := 1; i <= thread; i++ {
		go dispatcher.Dispatch(make(chan struct{}))
	}
	// delay period, maybe should push to kafka and safe another place.
	go dispatcher.recorder.startRecord()
}

// open 5 thread to handler events && message convert
func (d *_dispatcher) Dispatch(signal chan struct{}) {
	for {
		select {
		case scheme := <-d.schemeC:
			switch scheme.(type) {
			case *Event:
				event := scheme.(*Event)
				topic := event.GetTopic()
				broker, ok := gCache.GetNamespace(GetWorkerNamespace(topic))
				if !ok {
					logx.L().Warnf("namespace not found %v", topic)
					continue
				}
				_ = d.publisher.Publish(topic, broker.GetSelector().Select(), event.toReader())
			case *message, *Msg:
			}

		case <-semaphore:
			// receive a message notify signal than yield one to convert.
			d.lock.Lock()
			msg := heap.Pop(d.msgQ).(*message)
			d.lock.Unlock()

			err := d.convert(msg)
			if err != nil {
				logx.L().Errorf("_dispatcher asyncChan err %v", err)
				continue
			}

			runtime.Gosched()

		case <-d.ctx.Done():
			logx.L().Warnf("Dispatch thread shut down")
			runtime.Goexit()

		case <-signal:
			logx.L().Infof("Dispatch thread close automatically")
			runtime.Goexit()
		}
	}
}

func (d *_dispatcher) dispatchMany(msg []*message) error {
	for _, m := range msg {
		if m == nil {
			logx.L().Debugf("dispatchMany msg is nil")
			continue
		}
		err := d.convert(m)
		if err != nil {
			// todo: republic
			// nsq.ErrProtocol=E_PUB_FAILED PUB failed exiting
			logx.L().Errorf("dispatchMany %v", err)
			return err
		}
	}
	return nil
}

// illuminate: asyncChan and defaultHandlerMsg will calling this function.
// you can set breakpoint debug here for the panorama of the core schedule process.
func (d *_dispatcher) convert(msg *message) (err error) {
	if msg == nil {
		logx.L().Debugf("receive nil message in convert")
		return
	}
	switch msg.Meta.getMessageType() {
	// illuminate: ToTrigger will generate one toSlave message then convert to imputes.
	// ToServer message capsules with input and output, which reassembly in the db queue.
	case ToServer, ToTrigger, ToRet:
		// fill the context, only need pay attention to the ToServer/ToTrigger
		// before Execute, context provide abilities to executing sentence
		// to get all the node's input or output, this extra optional field may
		// enhance the data integrity, but it increase slot expenditure meanwhile.
		// so, it would be considerate to using gzip to compress context body in
		// the future.

		// now use the ttl to handler this is moderate(at 07.20). when a key
		// is never used just delete from context store. Output may not be map struct,
		// map[_internal_data:[12]] means that the data was sealed by a
		// map which key is _internal_data pointed.
		msg.Context.setSlaveOutput(msg.Task.getNodeCode(), msg.Task.getOutput())

		// question: it's valuable to pub/sub ToServer/ToTrigger message again?
		// to be honest, using queue to cut down the concurrent eps is valuable.
		// nevertheless, it seems that ToSlaver message should not keep in record.
		// sequence record bugs: defaultHandlerMsg will convert a message again,
		// ToSlave was the only type to be considered, but when we introduce ret
		// they are seems record, so the sequence overlapped cause the all
		// going through the case.
		err = d.handler(msg)
		if err != nil {
			msg.setExceptionError(err)
			logx.L().Errorf("_dispatcher internal error, %v", err)
		}

		// msg's messageType may changed when transform pointer to defaultHandlerMsg function.
		// which ToInternalError will be catch in the context of executing & calculating js edge.
		goto RECORD

	// illuminate: record && send event
	case ToException:
		// not return here, goto code will execute finally. so there is no
		// need to calling goto here.
		goto RECORD

	// Illuminate: ToSlave should not be record in database job.
	// meanwhile, only ToSlave need convert as a message express in the mq to
	// decline the eps, ToServer/ToException will pumped into local message channel in the db'job.
	// to persistence the record synchronously.
	// we need return message Publish before some nodes ToServer has already done.
	// so check the dependence is very basement. hook in Dispatch is wiser than any
	// places.
	case ToSlave:
		topic := msg.Meta.getTopic()
		if topic == "" {
			return fmt.Errorf("toSlave message topic is blank")
		}
		// not connected maybe
		ns, ok := gCache.GetNamespace(GetWorkerNamespace(topic))
		if !ok {
			return fmt.Errorf("topic %v not exist", GetWorkerNamespace(topic))
		}

		// in order to keep slave message in service palm, (e.g. ack check)
		// PutAck keep msg into exchange cache, tt not only reduces message size,
		// but also serialization time.
		exchange.PutAck(msg)
		gen := msg.genMsg()
		gen.Sign(Ecc)
		return d.publisher.Publish(topic, ns.GetSelector().Select(), gen.toReader())

	default:
		logx.L().Errorf("_dispatcher convert not implement yet: %v", msg.Meta.getMessageType())
		goto RECORD
	}

RECORD:
	switch msg.Meta.getMessageType() {
	case ToTrigger:
		// trigger sequence: trigger->slaver->server finally.
		var extra = Extra{}
		extra.Input = msg.Task.getInput()
		b, _ := json.Marshal(extra)
		d.recorder.executionChan <- &models.Execution{
			TraceId:    cast.ToString(msg.Meta.getTrace()),
			Sequence:   cast.ToInt(msg.Meta.getSequence()),
			NodeCode:   msg.Task.getNodeCode(),
			Status:     string(msg.Task.getStatus()),
			PlaybookId: msg.Task.getPlaybookId(),
			Extra:      string(b),
			Timestamp:  msg.Meta.getTimestamp(),
			Domain:     msg.Meta.getDomain(),
			SnapshotId: msg.Task.getSnapshotId(),
			Chain:      msg.Context.getChain(),
		}

	case ToServer, ToException, ToRet:
		var extra = Extra{}
		// getInput get the output from the worker handler.
		// getOutput get the input from the parentMessage.
		extra.Output = msg.Task.getOutput()
		extra.Input = msg.Task.getInput()

		var status = msg.Task.getStatus()
		// hook keeps all the context of message
		if msg.Task.isHook() {
			extra.Detail = msg
			status = Hooked
		}

		if msg.Meta.getMessageType() == ToException {
			extra.Exception = msg.Context.getException().GetContent()
		}

		b, _ := json.Marshal(extra)
		d.recorder.executionChan <- &models.Execution{
			TraceId:    cast.ToString(msg.Meta.getTrace()),
			Sequence:   cast.ToInt(msg.Meta.getSequence()),
			NodeCode:   msg.Task.getNodeCode(),
			PlaybookId: msg.Task.getPlaybookId(),
			Timestamp:  msg.Meta.getTimestamp(),
			Domain:     msg.Meta.getDomain(),
			SnapshotId: msg.Task.getSnapshotId(),
			Chain:      msg.Context.getChain(),
			Status:     string(status),
			Extra:      string(b),
		}

		// if exception or error in the playbook node executions environment
		if msg.Meta.getMessageType() == ToException {
			for msg.Context.getStack().size() > 0 {
				bury, retCode := msg.Context.getStack().pop()

				// when A call B, A and B has references, when B internal error happens, abort.
				// then delete slot key with domain, but A slot key would not delete, so slot leak
				// there is only pop when exception happens in this code, so delete here.
				// for further consideration, this block should be recode
				if pb, ok := gCache.GetPlaybook(bury.Task.getPlaybookId()); ok && pb.hasMark() {
					logx.L().Debugf("delete slot in pop stack")
					bury.delTrace()
				}

				if bury == nil {
					logx.L().Errorf("manipulate message maybe")
					continue
				}
				var extra = Extra{}
				extra.Input = bury.Task.getInput()
				extra.Exception = msg.Context.getException().GetContent()
				bs, _ := json.Marshal(extra)

				bury.Meta.addSequence()
				bury.Context.recordChain(retCode)

				d.recorder.executionChan <- &models.Execution{
					TraceId:    cast.ToString(msg.Meta.getTrace()),
					Sequence:   cast.ToInt(bury.Meta.getSequence()),
					NodeCode:   retCode,
					PlaybookId: bury.Task.getPlaybookId(),
					Timestamp:  bury.Meta.getTimestamp(),
					Domain:     bury.Meta.getDomain(),
					SnapshotId: bury.Task.getSnapshotId(),
					Chain:      bury.Context.getChain(),
					Extra:      string(bs),
					Status:     string(Fail),
				}
			}
		}
	default:
		logx.L().Errorf("not implement message'type: %v", msg.Meta.getMessageType())
	}
	return
}

// ----------------------------------------- Internal Queue

// non-thread-safe mq
type heapQ interface {
	heap.Interface
	Debug() string
}

type (
	priorityMq []*message
	fifoMq     chan *message
)

func (t *priorityMq) Len() int {
	return len(*t)
}

func (t *priorityMq) Less(i, j int) bool {
	// hook your strategy here
	if (*t)[i].Meta.Sync {
		return true
	}
	return (*t)[i].Meta.getMessageType() > (*t)[j].Meta.getMessageType()
}

func (t *priorityMq) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

func (t *priorityMq) Push(x interface{}) {
	*t = append(*t, x.(*message))
}

func (t *priorityMq) Pop() interface{} {
	n := len(*t)
	x := (*t)[n-1]
	*t = (*t)[:n-1]
	return x
}

// observer the internal heap status.
func (t *priorityMq) Debug() string {
	var count, other int
	for _, i := range *t {
		if i.Meta.getMessageType() == ToTrigger {
			count++
		} else {
			other++
		}
	}
	return fmt.Sprintf("ToTrigger count is: %d, others:%d", count, other)
}

func (t *fifoMq) Len() int {
	return len(*t)
}

func (t *fifoMq) Less(i, j int) bool {
	return false
}

func (t *fifoMq) Swap(i, j int) {

}

func (t *fifoMq) Push(x interface{}) {
	*t <- x.(*message)
}

func (t *fifoMq) Pop() interface{} {
	x := <-(*t)
	return x
}

func (t *fifoMq) Debug() string {

	return ""
}
