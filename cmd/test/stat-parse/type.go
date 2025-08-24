// A package of useful Go Kafka types
package main

/*
	ref - https://docs.confluent.io/platform/current/clients/librdkafka/html/md_STATISTICS.html

curl -I https://docs.confluent.io/platform/current/clients/librdkafka/html/md_STATISTICS.html
HTTP/2 200
content-type: text/html
content-length: 68898
date: Sat, 23 Aug 2025 02:50:44 GMT
last-modified: Fri, 22 Aug 2025 16:37:56 GMT
x-amz-version-id: MGISQ_dAnhHUn.cuhHhH0cpnRRJtIemJ
etag: "b8cf8e282117ef2f8eac90578959386d"
server: AmazonS3
x-cache: Hit from cloudfront
via: 1.1 b0921181e973f37be0dbce2713f6360e.cloudfront.net (CloudFront)
x-amz-cf-pop: IAH50-C2
x-amz-cf-id: WLrOGK0TUm3Noh1MjCqwYm4QFyzo-HjZbspyimIP2erINxkWpYwsXA==
age: 15
vary: Origin

Window stats
Rolling window statistics. The values are in microseconds unless otherwise stated.

Used for

type KafkaBroker
	IntLatency
	OutBufLatency
	Throttle

type KafkaTopic
	BatchSize
	BatchCnt

Key		Value	        Description
---		-----		----------------
min		int64		Smallest value
max		int64		Largest value
avg		int64		Average value
sum		int64		Sum of values
cnt		int64		Number of values sampled
stddev		int64		Standard deviation (based on histogram)
hdrsize		int64		Memory size of Hdr Histogram
p50		int64		50th percentile
p75		int64		75th percentile
p90		int64		90th percentile
p95		int64		95th percentile
p99		int64		99th percentile
p99_99		int64		99.99th percentile
outofrange	int64		Values skipped due to out of histogram range

*/

type LibRdKafkaStats struct {
	Name             string                 `json:"name,omitempty"`               // Handle instance name
	ClientId         string                 `json:"client_id,omitempty"`          // The configured (or default) client.id
	Type             string                 `json:"type,omitempty"`               // Instance type (producer or consumer)
	Ts               int64                  `json:"ts,omitempty"`                 // counter - librdkafka's internal monotonic clock (microseconds)
	Time             int64                  `json:"time,omitempty"`               // counter - Wall clock time in seconds since the epoch
	Age              int64                  `json:"age,omitempty"`                // counter - Time since this client instance was created (microseconds)
	ReplyQ           int64                  `json:"replyq,omitempty"`             // gauge - Number of ops (callbacks, events, etc) waiting in queue for application to serve with rd_kafka_poll()
	MsgCnt           int64                  `json:"msg_cnt,omitempty"`            // gauge - Current number of messages in producer queues
	MsgSize          int64                  `json:"msg_size,omitempty"`           // gauge - Current total size of messages in producer queues
	MsgMax           int64                  `json:"msg_max,omitempty"`            // counter - Threshold: maximum number of messages allowed allowed on the producer queues
	MsgSizeMax       int64                  `json:"msg_size_max,omitempty"`       // counter - Threshold: maximum total size of messages allowed on the producer queues
	Tx               int64                  `json:"tx,omitempty"`                 // counter - Total number of requests sent to Kafka brokers
	TxBytes          int64                  `json:"tx_bytes,omitempty"`           // counter - Total number of bytes transmitted to Kafka brokers
	Rx               int64                  `json:"rx,omitempty"`                 // counter - Total number of responses received from Kafka brokers
	RxBytes          int64                  `json:"rx_bytes,omitempty"`           // counter - Total number of bytes received from Kafka brokers
	TxMsgs           int64                  `json:"txmsgs,omitempty"`             // counter - Total number of messages transmitted (produced) to Kafka brokers
	TxMsgBytes       int64                  `json:"txmsg_bytes,omitempty"`        // counter - Total number of message bytes (including framing, such as per-Message framing and MessageSet/batch framing) transmitted to Kafka brokers
	RxMsgs           int64                  `json:"rxmsgs,omitempty"`             // counter - Total number of messages consumed, not including ignored messages (due to offset, etc), from Kafka brokers.
	RxMsgBytes       int64                  `json:"rsmxg_bytes,omitempty"`        // counter - Total number of message bytes (including framing) received from Kafka brokers
	SimpleCnt        int64                  `json:"simple_cnt,omitempty"`         // gauge - Internal tracking of legacy vs new consumer API state
	MetaDataCacheCnt int64                  `json:"metadata_cache_cnt,omitempty"` // gauge - Number of topics in the metadata cache
	Brokers          map[string]KafkaBroker `json:"brokers,omitempty"`            // Object list of brokers, key is broker name, value is object. See KafkaBrokers below
	Topics           map[string]KafkaTopic  `json:"topics,omitempty"`             // Object of topics, key is topic name, value is object. See KafkaTopics below
	Cgrp             KafkaCgrp              `json:"cgrp,omitempty"`               // Consumer group metrics. See KafkaCgrp below
	Eos              KafkaEos               `json:"eos,omitempty"`                // EOS / Idempotent producer state and metrics. See KafkaEos below
}

// Per broker metrics
type KafkaBroker struct {
	Name           string                `json:"name,omitempty"`             // Broker hostname, port and broker id e.g., "example.com:9092/13"
	NodeId         int64                 `json:"nodeid,omitempty"`           // Broker id (-1 for bootstraps)
	NodeName       string                `json:"nodename,omitempty"`         // Broker hostname e.g., "example.com:9092"
	Source         string                `json:"source,omitempty"`           // Broker source  (learned, configured, internal, logical)
	State          string                `json:"state,omitempty"`            // Broker state (INIT, DOWN, CONNECT, AUTH, APIVERSION_QUERY, AUTH_HANDSHAKE, UP, UPDATE)
	StateAge       int64                 `json:"stateage,omitempty"`         // gauge - Time since last broker state change (microseconds)
	OutBufCnt      int64                 `json:"outbuf_cnt,omitempty"`       // gauge - Number of requests awaiting transmission to broker
	OutBufMsgCnt   int64                 `json:"outbuf_msg_cnt,omitempty"`   // gauge - Number of messages awaiting transmission to broker
	WaitRespCnt    int64                 `json:"waitresp_cnt,omitempty"`     // gauge - Number of requests in-flight to broker awaiting response
	WaitRespMsgCnt int64                 `json:"waitresp_msg_cnt,omitempty"` // gauge - Number of messages in-flight to broker awaiting response
	Tx             int64                 `json:"tx,omitempty"`               // counter - Total number of requests sent
	TxBytes        int64                 `json:"txbytes,omitempty"`          // counter - Total number of bytes sent
	TxErrs         int64                 `json:"txerrs,omitempty"`           // counter - Total number of transmission errors
	TxRetries      int64                 `json:"txretries,omitempty"`        // counter - Total number of request retries
	TxIdle         int64                 `json:"txidle,omitempty"`           // counter - Microseconds since last socket send (or -1 if no sends yet for current connection).
	ReqTimeouts    int64                 `json:"req_timeouts,omitempty"`     // counter - Total number of requests timed out
	Rx             int64                 `json:"rx,omitempty"`               // counter - Total number of responses received
	RxBytes        int64                 `json:"rxbytes,omitempty"`          // counter - Total number of bytes received
	RxErrs         int64                 `json:"rxerrs,omitempty"`           // counter - Total number of receive errors
	RxCorridErrs   int64                 `json:"rxcorriderrs,omitempty"`     // counter - Total number of unmatched correlation ids in response (typically for timed out requests)
	RxPartial      int64                 `json:"rxpartial,omitempty"`        // counter - Total number of partial MessageSets received. The broker may return partial responses if the full MessageSet could not fit in the remaining Fetch response size.
	RxIdle         int64                 `json:"rxidle,omitempty"`           // counter - Microseconds since last socket receive (or -1 if no receives yet for current connection).
	Req            map[string]int64      `json:"req,omitempty"`              // Request type counters. Object key is the request name, value is the number of requests sent. See BrkrReqCounters
	ZBufGrow       int64                 `json:"zbuf_grow,omitempty"`        // counter - Total number of decompression buffer size increases
	BufGrow        int64                 `json:"buf_grow,omitempty"`         // counter - Total number of buffer size increases (deprecated, unused)
	WakeUps        int64                 `json:"wakeups,omitempty"`          // counter - Broker thread poll loop wakeups
	Connects       int64                 `json:"connects,omitempty"`         // counter - Number of connection attempts, including successful and failed, and name resolution failures.
	DisConnects    int64                 `json:"disconnects,omitempty"`      // counter - Number of disconnects (triggered by broker, network, load-balancer, etc.).
	IntLatency     map[string]int64      `json:"int_latency,omitempty"`      // Internal producer queue latency in microseconds. See Window stats below
	OutBufLatency  map[string]int64      `json:"outbuf_latency,omitempty"`   // Internal request queue latency in microseconds. This is the time between a request is enqueued on the transmit (outbuf) queue and the time the request is written to the TCP socket. Additional buffering and latency may be incurred by the TCP stack and network. See Window stats below
	RTT            map[string]int64      `json:"rtt,omitempty"`              // Broker latency / round-trip time in microseconds. See Window stats below
	Throttle       map[string]int64      `json:"throttle,omitempty"`         // Broker throttling time in milliseconds. See Window stats below
	TopPars        map[string]BrokerPars `json:"toppars,omitempty"`          // Partitions handled by this broker handle. Key is "topic-partition". See brokers.toppars below
}

type KafkaTopic struct {
	Topic       string                    `json:"topic,omitempty"`        // Topic name
	Age         int64                     `json:"age,omitempty"`          // Age of client's topic object (milliseconds)
	MetaDataAge int64                     `json:"metadata_age,omitempty"` // Age of metadata from broker for this topic (milliseconds)
	BatchSize   map[string]int64          `json:"batchsize,omitempty"`    // Batch sizes in bytes. See *Window stats*·
	BatchCnt    map[string]int64          `json:"batchcnt,omitempty"`     // Batch message counts. See *Window stats*·
	Partitions  map[string]KafkaPartition `json:"partitions,omitempty"`   // Topic partitions
}

type KafkaPartition struct {
	Partition         int64  `json:"partition,omitempty"`           // Partition Id (-1 for internal UA/UnAssigned partition)
	Broker            int64  `json:"broker,omitempty"`              // The id of the broker that messages are currently being fetched from
	Leader            int64  `json:"leader,omitempty"`              // Current leader broker id
	Desired           bool   `json:"desired,omitempty"`             // Partition is explicitly desired by application
	Unknown           bool   `json:"unknown,omitempty"`             // Partition not seen in topic metadata from broker
	MsgQCnt           int64  `json:"msgq_cnt,omitempty"`            // gauge - Number of messages waiting to be produced in first-level queue
	MsgQBytes         int64  `json:"msgq_bytes,omitempty"`          // gauge - Number of bytes in msgq_cnt
	XmitMsgQCnt       int64  `json:"xmit_msgq_cnt,omitempty"`       // gauge - Number of messages ready to be produced in transmit queue
	XmitMsgBytes      int64  `json:"xmit_msgq_bytes,omitempty"`     // gauge - Number of bytes in xmit_msgq
	FetchQCnt         int64  `json:"fetchq_cnt,omitempty"`          // gauge - Number of pre-fetched messages in fetch queue
	FetchQSize        int64  `json:"fetchq_size,omitempty"`         // gauge - Bytes in fetchq
	FetchState        string `json:"fetch_state,omitempty"`         // Consumer fetch state for this partition (none, stopping, stopped, offset-query, offset-wait, active)
	QueryOffset       int64  `json:"query_offset,omitempty"`        // gauge - Current/Last logical offset query
	NextOffset        int64  `json:"next_offset,omitempty"`         // gauge - Next offset to fetch
	AppOffset         int64  `json:"app_offset,omitempty"`          // gauge - Offset of last message passed to application + 1
	StoredOffset      int64  `json:"stored_offset,omitempty"`       // gauge - Offset to be committed
	StoredLeaderEpoch int64  `json:"stored_leader_epoch,omitempty"` // Partition leader epoch of stored offset
	// Note: there is a typo in the stats.  Unsure which will stay.  Take both for now
	CommitedOffset       int64 `json:"commited_offset,omitempty"`        // gauge - Last committed offset
	CommittedOffset      int64 `json:"committed_offset,omitempty"`       // gauge - Last committed offset
	CommittedLeaderEpoch int64 `json:"committed_leader_epoch,omitempty"` // Partition leader epoch of committed offset
	EofOffset            int64 `json:"eof_offset,omitempty"`             // gauge - Last PARTITION_EOF signaled offset
	LoOffset             int64 `json:"lo_offset,omitempty"`              // gauge - Partition's low watermark offset on broker
	HiOffset             int64 `json:"hi_offset,omitempty"`              // gauge - Partition's high watermark offset on broker
	LsOffset             int64 `json:"ls_offset,omitempty"`              // gauge - Partition's last stable offset on broker, or same as hi_offset if broker version is less than 0.11.0.0.
	ConsumerLag          int64 `json:"consumer_lag,omitempty"`           // gauge - Difference between (hi_offset or ls_offset) and committed_offset). hi_offset is used when isolation.level=read_uncommitted, otherwise ls_offset.
	ConsumerLagStored    int64 `json:"consumer_lag_stored,omitempty"`    // gauge - Difference between (hi_offset or ls_offset) and stored_offset. See consumer_lag and stored_offset.
	LeaderEpoch          int64 `json:"leader_epoch,omitempty"`           // Last known partition leader epoch, or -1 if unknown.
	TxMsgs               int64 `json:"txmsgs,omitempty"`                 // counter - Total number of messages transmitted (produced)
	TxBytes              int64 `json:"txbytes,omitempty"`                // counter - Total number of bytes transmitted for txmsgs
	RxMsgs               int64 `json:"rxmsgs,omitempty"`                 // counter - Total number of messages consumed, not including ignored messages (due to offset, etc).
	RxBytes              int64 `json:"rxbytes,omitempty"`                // counter - Total number of bytes received for rxmsgs
	Msgs                 int64 `json:"msgs,omitempty"`                   // counter - Total number of messages received (consumer, same as rxmsgs), or total number of messages produced (possibly not yet transmitted) (producer)
	RxVerDrops           int64 `json:"rx_ver_drops,omitempty"`           // counter? - Dropped outdated messages
	MsgsInflight         int64 `json:"msgs_inflight,omitempty"`          // gauge - Current number of messages in-flight to/from broker
	NextAckSeq           int64 `json:"next_ack_seq,omitempty"`           // gauge - Next expected acked sequence (idempotent producer)
	NextErrSeq           int64 `json:"next_err_seq,omitempty"`           // gauge - Next expected errored sequence (idempotent producer)
	AckedMsgId           int64 `json:"acked_msgid,omitempty"`            // Last acked internal message id (idempotent producer)
}

// Topic partitions assigned to broker
type BrokerPars struct {
	Topic     string `json:"topic,omitempty"`     // topic name
	Partition int64  `json:"partition,omitempty"` // partition id
}

type KafkaCgrp struct {
	State           string `json:"state,omitempty"`            // Local consumer group handler's state.
	StateAge        int64  `json:"stateage,omitempty"`         // gauge - Time elapsed since last state change (milliseconds)
	JoinState       string `json:"join_state,omitempty"`       // Local consumer group handler's join state
	RebalanceAge    int64  `json:"rebalance_age,omitempty"`    // gauge - Time elapsed since last rebalance (assign or revoke) (milliseconds)
	RebalanceCnt    int64  `json:"rebalance_cnt,omitempty"`    // counter - Total number of rebalances (assign or revoke)
	RebalanceReason string `json:"rebalance_reason,omitempty"` // Last rebalance reason, or empty string
	AssignmentSize  int64  `json:"assignment_size,omitempty"`  // gauge - Current assignment's partition count
}

type KafkaEos struct {
	IdempState    string `json:"idemp_state,omitempty"`    // Current idempotent producer id state.
	IdempStateAge int64  `json:"idemp_stateage,omitempty"` // gauge - Time elapsed since last idemp_state change (milliseconds)
	TxnState      string `json:"txn_state,omitempty"`      // Current transactional producer state
	TxnStateAge   int64  `json:"txn_stateage,omitempty"`   // counter - Time elapsed since last txn_state change (milliseconds)
	TxnMayEnq     bool   `json:"txn_may_enq,omitempty"`    // Transactional state allows enqueuing (producing) new messages
	ProducerId    int64  `json:"producer_id,omitempty"`    // gauge - The currently assigned Producer ID (or -1)
	ProducerEpoch int64  `json:"producer_epoch,omitempty"` // gauge - The current epoch (or -1)
	EpochCnt      int64  `json:"epoch_cnt,omitempty"`      // The number of Producer ID assignments since start
}
