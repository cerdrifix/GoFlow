
create temporary table _instances (
  inst_id uuid,
  map_id uuid
);

create temporary table _nodes (
    inst_id uuid,
    node_name varchar(255),
    node_description varchar(255)
);

insert into _instances (inst_id, map_id)
select      id, map_id
from        instances;

insert into _nodes (inst_id, node_name, node_description)
select i.inst_id,
       json_array_elements(m.data -> 'nodes') ->> 'name',
       json_array_elements(m.data -> 'nodes') ->> 'description'
from        maps m
inner join  _instances i on m.id = i.map_id;

select i.id as instance_id,
       n.node_description as current_node,
       i.start_date,
       s.creator_id,
       concat(u.surname, ' ', u.name) as creator_description
from       instances i
inner join states s on i.current_state = s.id
inner join _nodes n on n.inst_id = i.id and n.node_name = s.node_name
inner join users u on s.creator_id = u.username;

drop table if exists _instances;
drop table if exists _nodes;