# Tourist-itinerary
# Description of task

For building tourist routes, it is necessary to take into account all possible ways of moving between attractions. This is on foot, by public transport and by taxi. In this chapter, we will describe algorithms for constructing itinerary on public transport. The results of the algorithm should be as close as possible to the real ones. Therefore, the route will take into account the bus schedule and traffic jams on the road in real time.
The construction of this model will be implemented in several steps:
- Parsing data about public transport routes;
- Parsing data about traffic-jams in real time;
- Processing the received data to obtain general information about itinerary on public transport;
- Building direct route by public transport;
- Building routes with one transfer by public transport;
- Building routes with two transfers by public transport.

All these steps will help you build routes in accordance with the actual situation on the roads. In order to adapt the routes to comfortable conditions for tourists, the following chapters will consider modifications of this algorithm. 
The algorithm optimization priority goes in the following order:
- Minimize the number of transfers;
- Minimize the time spent on the trip.

The option with three or more transfers is not considered because in a large city this is quite a rare phenomenon, and it is also as uncomfortable as possible for tourists. Also, all the popular attractions of the city are usually within the reach of a maximum of two transfers. Otherwise, perhaps this attraction probably does not deserve a visit.


# Detailed description is presented in attached PDF file.
