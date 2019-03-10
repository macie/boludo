
reservation(Location) :- atom_string(Location, "reservation://").

reservation(Resource, Location) :-  % reservation://<resources>/
    reservation(Base_uri),
    atom_concat(Resource, "/", Resource_uri),
    atom_concat(Base_uri, Resource_uri, Location).

reservation_rooms(Location) :-  % reservation://rooms/
    reservation("rooms", Location).

reservation_rooms(Room, Location) :-  % reservation://rooms/<room_name>/
    reservation_rooms(Base_uri),
    atom_concat(Room, "/", Room_uri),
    atom_concat(Base_uri, Room_uri, Location).

reservation_owner(Location) :-  % reservation://rooms/owners/
    reservation_owner("", Location).

reservation_owner(Room, Location) :-  % reservation://rooms/<room_name>/owners/
    reservation_rooms(Room, Base_uri),
    atom_concat(Base_uri, "owners/", Location).

reservation_free(Location) :-  % reservation://rooms/free/
    reservation_free("", Location).

reservation_free(Room, Location) :-  % reservation://rooms/<room_name>/free/
    reservation_rooms(Room, Base_uri),
    atom_concat(Base_uri, "free/", Location).

% reservation://users
% reservation://users/<username>/
% reservation://users/<username>/reservations
% ...
