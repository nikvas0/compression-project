# compression-project

Запуск:

0) Запуск docker-compose build && docker-compose up

1) Создаем стрим curl -v localhost:8081/start_stream?path=test

2) В obs studio транслируем на сервер rtmp://localhost:1935 с key = /test

3) В vlc/ffplay можно открыть http://localhost:8082/hls/test-360p.m3u8 и http://localhost:8082/hls/test-720p.m3u8.
Они отличаются значениями разрешения и битрейта.


---------------------


Видео-стриминг в реальном времени.

[x] 4 балла - сервер принимает поток видео и ретранслирует его

[ ] +2 балла - клиент --- TODO: отдача из manager_service <video>...</video> для начала, затем где-то форма для создания своего стрима

[ ] +2 балла - CDN --- в последнюю очередь

[x] +3 балла - при ретрансляции возможно изменение битрейта --- сейчас у нас можно выбрать между 360p и 720p

[ ] +3 балла - real-time чат --- TODO: любым образом + html клиент к чату, выданный из manager_service


