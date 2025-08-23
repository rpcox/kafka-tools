## output

    > ./stat-parse
                    Name: rdkafka#consumer-1
                ClientID: rdkafka
                    Type: consumer
                      Ts: 1754362881471038 (2025-08-05 03:01:21 UTC)
                    Time: 1754362881 (2025-08-05 03:01:21 UTC)
                     Age: 5004758 (3m0s)
                  ReplyQ: 0
                  MsgCnt: 0
                 MsgSize: 0
                  MsgMax: 0
              MsgSizeMax: 0
                      Tx: 23
                 TxBytes: 1902
                      Rx: 20
                 RxBytes: 4437
                  TxMsgs: 0
             TxMsgsBytes: 0
                  RxMsgs: 0
             RxMsgsBytes: 0
               SimpleCnt: 0
        MetaDataCacheCnt: 1

                BROKERS:
    ++       smaug:9092/300
                  NodeId: 300
                   State: UP
              Req[Fetch]: 2
         IntLatency[min]: 0
      OutBufLatency[min]: 31
                RTT[avg]: 285789 (7.2s)
           Throttle[avg]: 0 (0s)

      BROKER PARTITIONS:
              test.topic: 0
              test.topic: 3

    ++        puff:9092/200
                  NodeId: 200
                   State: UP
              Req[Fetch]: 2
         IntLatency[min]: 0
      OutBufLatency[min]: 27
                RTT[avg]: 285710 (7.2s)
           Throttle[avg]: 0 (0s)

      BROKER PARTITIONS:
              test.topic: 2
              test.topic: 5

    ++       mushu:9092/100
                  NodeId: 100
                   State: UP
              Req[Fetch]: 2
         IntLatency[min]: 0
      OutBufLatency[min]: 37
                RTT[avg]: 348566 (10.8s)
           Throttle[avg]: 0 (0s)

      BROKER PARTITIONS:
              test.topic: 1
              test.topic: 4


                 TOPICS:
    ++           test.topic
      P-0 Consumer Lag: : 0
      P-1 Consumer Lag: : -1
      P-2 Consumer Lag: : 0
      P-3 Consumer Lag: : 0
      P-4 Consumer Lag: : 0
      P-5 Consumer Lag: : 0

    >
