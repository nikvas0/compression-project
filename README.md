# compression-project

Запуск:

0) Запуск docker-compose build && docker-compose up

1) Заходим на http://localhost:7191/home

2) Создаем стрим из интерфейса

3) В obs studio транслируем на сервер rtmp://localhost:1935 с key = /your_stream_name

4) В vlc/ffplay можно открыть http://localhost:8082/hls/test_360p.m3u8 и http://localhost:8082/hls/test_720p.m3u8.
Они отличаются значениями разрешения и битрейта.
Также теперь можно зайти через интерфейс на http://localhost:7191/streams и выбрать нужный стрим, после чего нажать play (браузеры запрещают autoplay без mute).



---------------------


Видео-стриминг в реальном времени.

- [x] 4 балла - сервер принимает поток видео и ретранслирует его

- [x] +2 балла - клиент --- TODO: где-то форма для создания своего стрима 

- [ ] +2 балла - CDN

- [x] +3 балла - при ретрансляции возможно изменение битрейта --- сейчас у нас можно выбрать между разрешением 360p и 720p и двумя битрейтами, которые им соответствуют

- [ ] +3 балла - real-time чат --- TODO: любым образом + html клиент к чату, выданный из manager_service


