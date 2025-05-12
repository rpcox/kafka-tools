## kpload


    > cat file
    {"first":"john","last":"smaith"}
    {"first":"jane","last":"doe"}

    > cat file | kpload -p props -t test.topic -keygen enum
    publish success: topic:test.topic in partition[5] @ offset 9
    publish success: topic:test.topic in partition[5] @ offset 10

    > cat file | kpload -p props -t test.topic -keygen epoch
    publish success: topic:test.topic in partition[0] @ offset 10
    publish success: topic:test.topic in partition[0] @ offset 11

    > cat file | ./kpload -p props -t test.topic -keygen uuid
    publish success: topic:test.topic in partition[5] @ offset 11
    publish success: topic:test.topic in partition[5] @ offset 12


    > kafka-console-consumer.sh --topic test.topic --bootstrap-server smaug:9092 --from-beginning --property print.key=true
    1	{"first":"john","last":"smaith"}
    1	{"first":"jane","last":"doe"}
    1747020727	{"first":"john","last":"smaith"}
    1747020727	{"first":"jane","last":"doe"}
    b8733bce-0372-497b-b245-e88e8d787c0b	{"first":"john","last":"smaith"}
    57b67ee6-327d-46a1-bb71-52f5f0118459	{"first":"jane","last":"doe"}
