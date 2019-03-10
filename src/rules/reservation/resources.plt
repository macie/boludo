:- begin_tests(resources).
:- include(resources).

test(reservation) :- reservation('reservation://').

test(reservation_rooms) :- reservation_rooms('reservation://rooms/').
test(reservation_rooms) :- reservation_rooms(11, 'reservation://rooms/11/').

test(reservation_owner) :- reservation_owner('reservation://rooms/owners/').
test(reservation_owner) :- reservation_owner(5, 'reservation://rooms/5/owners/').

test(reservation_free) :- reservation_free('reservation://rooms/free/').
test(reservation_free) :- reservation_free(3, 'reservation://rooms/3/free/').

:- end_tests(resources).
