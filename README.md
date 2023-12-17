# WBTech L0 project

NATS-streaming-server и экземпляр postgres развернуты в docker контейнерах.
Для их запуска
```bash
make run
```

Запуск stan-publisher и публикация тестовых данных
```bash
make producer
```

Запуск stan-subscriber
```bash
make consumer
```