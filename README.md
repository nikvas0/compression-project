# compression-project

Запуск:

0) Запуск docker-compose build && docker-compose up

1) Создаем стрим curl -v localhost:8081/start_stream?path=test

2) В obs studio транслируем на сервер rtmp://localhost:1935 с key = /test

3) В vlc/ffplay можно открыть http://localhost:8082/hls/test_360p.m3u8 и http://localhost:8082/hls/test_720p.m3u8.
Они отличаются значениями разрешения и битрейта.
Таже можно теперь открыть трансляции в браузере по ссылкам http://localhost:7191/video?path=test&r=720p и http://localhost:7191/video?path=test&r=360p.



---------------------


Видео-стриминг в реальном времени.

- [x] 4 балла - сервер принимает поток видео и ретранслирует его

- [ ] +2 балла - клиент --- TODO: отдача из manager_service \<video\>: нужно доделать mute/unmute + где-то форма для создания своего стрима

- [ ] +2 балла - CDN

- [x] +3 балла - при ретрансляции возможно изменение битрейта --- сейчас у нас можно выбрать между разрешением 360p и 720p и двумя битрейтами, которые им соответствуют

- [ ] +3 балла - real-time чат --- TODO: любым образом + html клиент к чату, выданный из manager_service


