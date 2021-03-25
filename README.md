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

- [x] +2 балла - клиент

- [x] +3 балла - при ретрансляции возможно изменение битрейта --- сейчас у нас можно выбрать между разрешением 360p и 720p и двумя битрейтами, которые им соответствуют

- [X] +3 балла - real-time чат

- [ ] +2 балла - CDN --- НЕТ



---------------------

![arch](arch.jpg)


Твич использует RTMP для отправки видео и HLS streams для трансляции.Будем делать все также, так как HLS поддерживается браузерами нативно(имеет поддержку nginx + просто работать с CDN), используем его.Сервис будет принимать RTMP поток, затем перекодировать через ffmpeg(несколько разных битрейтов сразу, чтобы его можно было выбирать) и раздавать трансляцию через nginx


manager service -- сервис, предоставляющий REST API для управлениястримами

video service -- сервис, преобразующий RTMP в HLS используя ffmpeg

chat service -- сервис чатов

### Трансляция видео

vides service запускает ретрансляцию rtmp потока в ffmpeg с разными разрешениями/битрейтами, который переделывает поток в HLS формат и складывает в файлы, раздающиеся через nginx.





