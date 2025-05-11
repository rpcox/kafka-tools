 ## kp1

    >  bin/kafka-topics.sh --bootstrap-server mushu:9092 --list --exclude-internal
    syslog-systemd
    test.topic
    > kp1 -p props -t test.topic
    publish success: topic:test.topic in partition[3] @ offset 0
    publish success: topic:test.topic in partition[3] @ offset 1
    publish success: topic:test.topic in partition[3] @ offset 2
    publish success: topic:test.topic in partition[3] @ offset 3
    publish success: topic:test.topic in partition[3] @ offset 4
    > kp1 -p props -t test.topic
    publish success: topic:test.topic in partition[3] @ offset 5
    publish success: topic:test.topic in partition[3] @ offset 6
    publish success: topic:test.topic in partition[3] @ offset 7
    publish success: topic:test.topic in partition[3] @ offset 8
    publish success: topic:test.topic in partition[3] @ offset 9
     > kafka-console-consumer.sh --topic test.topic --bootstrap-server smaug:9092 --from-beginning
    0:ByOagAheU6ximq2WK77HQzeO7ZjQQ9nNHmSQnXa1R3iLag3fJHvAjJRKMfTVU1ZA
    1:chdbVYxWKl4yNz9y1GslKGTe7Xy0aFjDTwkO3EHVw76E8l3tYte4BGcx14f8lICY
    2:FPyyuS0ioprBag7yfPZgvPhrpzIL99Rvj8S6xda71SQ74hnJbpBw6JBTN0eCZRbY
    3:YVB9HsfO1r2k90tWTaczC0mcOvhchPrzgKhiV1gdckeLgLEhuHDofuwmApHoVexQ
    4:A1ZjP1oeVAv2r3nVC9gWdHMpdqxK1FrQiDYD2wqNpsLd7bmyzFiB6wnnXC2hAvII
    0:masahXc8Ssl1kC8YOAPxdjK4hQXkGGw9C3Cy8KlEETmBcb8aa41aZN9a6LhlVygA
    1:xsbsUmwyTPeXYLvkYhM92EvNSRl4xuIXNxofGcxftmL5YuJy6vuFBnANHuGrQeXl
    2:h9NUYeVE9fJmrDky4N4tIWQuYRvhUU5RgMQVwsf5bxyx3XwChVrDh4ZZNCyLVlmF
    3:IIlfwynzD54yCabSBYchtFRfT7UJhCkhyFlK0qh3zVZpSk51Ri5Oswu11g26Pkt6
    4:JFxSp52m4sNOVRjqPG7QN4mK8f49QoAewrXVOPjkxy8ZKAhkbGEpqgceDHwRe6Pc
    ^C
    >
