CREATE TABLE `app` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_name` varchar(255) NOT NULL DEFAULT '' COMMENT 'app名称',
  `user` varchar(255) NOT NULL DEFAULT '' COMMENT '用户',
  `token` varchar(255) NOT NULL DEFAULT '' COMMENT '是否有权限消费nsq',
  `brokers` varchar(255) NOT NULL DEFAULT '' COMMENT 'brokers',
  `broker_type` varchar(255) NOT NULL DEFAULT '' COMMENT '消息队列类型',
  `eps` int(255) NOT NULL DEFAULT 1 COMMENT '限流',
  `group_id` varchar(11) NOT NULL DEFAULT '' COMMENT '用户字段组',
  `update_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '更新时间',
  `checked` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否通过核审',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '信息描述',
  `last_alive_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '上次心跳时间',
  `share` tinyint(1) NOT NULL DEFAULT 0 COMMENT '共享app',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_name` (`app_name`) USING BTREE,
  UNIQUE KEY `index_token` (`token`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1  COMMENT='App表用于标注一个命名空间';


CREATE TABLE `app_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '外键',
  `app_id` int(11) NOT NULL DEFAULT 0 COMMENT 'app外键',
  `checked` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否通过',
  `create_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='关系';

CREATE TABLE `execution_0` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';

CREATE TABLE `execution_1` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';

CREATE TABLE `execution_2` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';


CREATE TABLE `execution_3` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';

CREATE TABLE `execution_4` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';


CREATE TABLE `execution_5` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';


CREATE TABLE `execution_6` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';


CREATE TABLE `execution_7` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';


CREATE TABLE `execution_8` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';


CREATE TABLE `execution_9` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `trace_id` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路id',
  `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '节点code',
  `node_code` varchar(32) NOT NULL DEFAULT '' COMMENT '节点id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '所属域',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '剧本id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0 COMMENT '快照id',
  `extra` text COMMENT '额外信息',
  `timestamp` varchar(55) NOT NULL DEFAULT '' COMMENT '时间戳',
  `chain` varchar(255) NOT NULL DEFAULT '' COMMENT '执行链路',
  PRIMARY KEY (`id`),
  KEY `index_trace_id` (`trace_id`) USING BTREE,
  KEY `index_node_code` (`node_code`) USING BTREE,
  KEY `index_snapshot_id` (`snapshot_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';
CREATE TABLE `playbook` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_id` int(11) NOT NULL DEFAULT 0 COMMENT '对应的app id',
  `snapshot_id` int(11) NOT NULL DEFAULT 0  COMMENT '快照版本id，用于动态切换',
  `user` varchar(255) NOT NULL DEFAULT '' COMMENT '创建人',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '剧本名称',
  `enable` tinyint(255) NOT NULL DEFAULT 0 COMMENT '是否开启',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '剧本描述',
  `token` varchar(255) NOT NULL DEFAULT '' COMMENT '授权token，用于被远程调用时鉴权',
  `update_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行' ;

CREATE TABLE `snapshot` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `playbook_id` int(11) NOT NULL DEFAULT 0 COMMENT '对应的剧本id',
  `snapshot` text  COMMENT '剧本快照body体',
  `rawbody` text  COMMENT '前端传来的元数据',
  `checksum` varchar(255) NOT NULL DEFAULT '' COMMENT '校验和',
  `app_id` int(11) NOT NULL DEFAULT 0 COMMENT 'app',
  `snapname` varchar(255) NOT NULL DEFAULT '' COMMENT '快照名称',
  `update_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '更新时间',
  `user` varchar(255) NOT NULL default '' COMMENT '更新快照人',
  PRIMARY KEY (`id`),
  KEY `index_pb` (`playbook_id`) USING BTREE COMMENT '剧本id索引',
  KEY `index_app` (`app_id`) USING BTREE COMMENT 'appid索引',
  KEY `index_checksum` (`checksum`) USING BTREE COMMENT '校验和索引'
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='执行';

CREATE TABLE `tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '函数名',
  `configuration` varchar(255) NOT NULL DEFAULT '' COMMENT '函数的配置',
  `xrn` varchar(255) NOT NULL DEFAULT '' COMMENT '对应的xrn',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `app_id` int(11) NOT NULL default 0 COMMENT '对应的appid',
  `type` varchar(255) NOT NULL default '' COMMENT '任务类型',
  `update_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '更新时间',
  `input_example` varchar(255) NOT NULL DEFAULT '' COMMENT '输入样例',
  `output_example` varchar(255) NOT NULL DEFAULT '' COMMENT '输出样例',
  `user` varchar(255) NOT NULL default '' COMMENT '创建者',
  PRIMARY KEY (`id`),
  KEY `aoo_index` (`app_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1  COMMENT='任务表';

CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `account` varchar(255) NOT NULL DEFAULT '' COMMENT '用户',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱',
  `workcode` varchar(255) NOT NULL DEFAULT '' COMMENT '工号',
  `dept_id` varchar(255) NOT NULL DEFAULT '' COMMENT '部门',
  `dept_name` varchar(255) NOT NULL DEFAULT '' COMMENT '部门',
  `dept_full_name` varchar(255) NOT NULL DEFAULT '' COMMENT '部门',
  `email` varchar(255) NOT NULL DEFAULT '' COMMENT 'email',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `create_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '创建时间',
  `super_admin` tinyint(1) NOT NULL DEFAULT 0 COMMENT '超级管理员',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='用户表';
