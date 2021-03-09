# compression-project

Запуск:

0) Запуск docker-compose build && docker-compose up

1) Создаем стрим curl -v localhost:8081/start_stream?path=test

2) В obs studio транслируем на сервер rtmp://localhost:1935 с key = /test

3) В vlc/ffplay можно открыть http://localhost:8082/hls/test-360p.m3u8 и http://localhost:8082/hls/test-720p.m3u8.
Они отличаются значениями разрешения и битрейта.

