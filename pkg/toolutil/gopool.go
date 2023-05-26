// Author: huaxr
// Time: 2022/6/17 2:38 下午
// Git: huaxr

package toolutil

type poolChain struct {
	elem chan struct{}
	next *poolChain
	size int32
}

func (p *poolChain) add() {
	tmp := p
	k := poolChain{
		elem: make(chan struct{}),
		next: tmp,
	}
	p.size += 1
	p = &k
}

func (p *poolChain) del() {
	tmp := p.next
	p = tmp
}
