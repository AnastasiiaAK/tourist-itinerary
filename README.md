# Tourist-itinerary
Основной алгорит реализован на языку Golang (код можно посмотреть в папке Golang). Обработка данных о пробках, аналитический обзор и алгоритмы машинного обучения реализованы  на Python. Подробное описание работы можно найти в приложенном файле pdf.
1. На вход программе передается 2 достопримечательности, между которыми турист хочет проложить оптимальный маршрут
(также турист может выбрать коэффициент предпочтения пешей прогулки или поездке на транспорте). Если коэффициент большой, то турист проходит максимально возможную в разумных пределах часть маршрута пешком. Если коэффициент маленький, то турист как можно больше времени проезжает в транспорте.

Варианты преодоления пути:
1. Пеший маршрут.
Данные о дорогах для построения пешего маршрута были скачаны с Overpass API. Маршрут строится с помощью алгоритма Дейкстры по пешеходным дорогам.

2. Маршрут на наземном общественном транспорте
Данные о маршрута общественного транспорта взяты с сайта общественного транспорта Санкт-Петербурга.
Данные о пробках получены с помощью данных о реальном местоположении конкретного транспорта в текущий момент времени (данные скачивались каждые 3 минуты в течении недели). Далее полученные данные обрабатывались. Если автобус был в пределах 100 метрах от остановки - то текущее время являлось временем прихода на конкретную остановку для конкретного транспорта. Далее получился датасет на конкретный интервал в 1 час, который содержит время прихода на конкретную остановку для конкретного транспорта.  Могло так получится, что время на каких-то остановках неизвестно, но считалось как расстояние деленное на скорость (считалась как расстояние деленное на время между ближайшими остановками с ближайшим временем) и с помощью алгоритмов машинного обучения (где на вход подавались следующие параметры - долгота, широта, расстяние между оставноками, расстояние до центра и направление движения (в сторону центра или из него), время движения, день недели и сезон уже учитывается, потому что модель обучается на данных только в конкретном интервале в один час).
Далее с на полученных данных строились маршруты общественного транспорта.
Определяется  время начала движения, в зависимости от него подбираются конкретные файлы с данными о загруженности дорог в текущий момент времени и строились маршруты (прямой,с одной пересадкой, с двумя пересадками).
Данные об остановках хранятся с помощью k-d tree структуры, что позволяет за довольно быстрое время подбирать набор ближайших остановок к точке начала и конца.
Подробные алгоритмы описаны в pdf файле и реализованы в папке golang.

3. Маршруты на такси.
Строились с помощью полученных на предыдущем шаге данных о загруженности дорог в конкретном промежуток времени. Данные подавались как веса в алгоритм Дейкстры. Данные о дорогах получены спарсены с Overpass API.


Данный алгоритм предназначен для построения маршрутов туристами. В моей работе он работал с фреймворком для построения маршрутов по достопримечательностям (когда на вход дается свободное время и набор достопримечательностей и задачей фреймворка является подобрать оптимальный набор достопримечательной с наибольшей полезностью и с ограниченным временем)




Необходимые данные:

Точка старта, точка финиша и планируемо время старта (если время не задано, то маршрут строится на текущее время) и коэффициент предпочтения пешего маршрута (если не задан, то коэффициент равен 1 и время маршрута пешком может достигать 20 минут от точки старта и 20 минут до точки финиша). 
Программа выдает лучший маршрут с учетом требований и наименьшим временем в пути.



# Description of task

For building tourist routes, it is necessary to take into account all possible ways of moving between attractions. This is on foot, by public transport and by taxi. In this chapter, we will describe algorithms for constructing itinerary on public transport. The results of the algorithm should be as close as possible to the real ones. Therefore, the route will take into account the bus schedule and traffic jams on the road in real time.
The construction of this model will be implemented in several steps:
- Parsing data about public transport routes;
- Parsing data about traffic-jams in real time;
- Processing the received data to obtain general information about itinerary on public transport;
- Ananlyze receive data and predict some unknown data;
- Building direct route by public transport;
- Building routes with one transfer by public transport;
- Building routes with two transfers by public transport.

All these steps will help you build routes in accordance with the actual situation on the roads. In order to adapt the routes to comfortable conditions for tourists, the following chapters will consider modifications of this algorithm. 
The algorithm optimization priority goes in the following order:
- Minimize the number of transfers;
- Minimize the time spent on the trip.

The option with three or more transfers is not considered because in a large city this is quite a rare phenomenon, and it is also as uncomfortable as possible for tourists. Also, all the popular attractions of the city are usually within the reach of a maximum of two transfers. Otherwise, perhaps this attraction probably does not deserve a visit.


# Detailed description is presented in attached PDF file.
