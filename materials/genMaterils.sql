CREATE TABLE materials
(
id SERIAL PRIMARY KEY,
pid VARCHAR(255) NOT NULL,
material_code VARCHAR(255) NOT NULL,
material_group VARCHAR(255) NOT NULL,
discription VARCHAR(255),
material_name VARCHAR(255) NOT NULL,
material_unit VARCHAR(20),
base_price decimal,
gst INT,
created_at TIMESTAMP NOT NULL,
created_by VARCHAR(255) NOT NULL,
last_updated_by VARCHAR(255) NOT NULL,
last_updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX material_code_material_group_idxs ON materials (material_code,material_group);
CREATE TABLE trade(
trade_id serial PRIMARY KEY,
trade_name varchar(100),
trade_description varchar(100),
is_material SMALLINT DEFAULT 0 NOT NULL,
is_organization SMALLINT DEFAULT 0 NOT NULL,
status SMALLINT DEFAULT 0 NOT NULL,
created_by VARCHAR(50) NOT NULL,
created_at TIMESTAMP NOT NULL,
last_updated_by VARCHAR(50) NOT NULL,
last_updated_at TIMESTAMP NOT NULL
);

INSERT INTO public.trade(
     trade_name, is_material, is_organization, status, created_by, created_at, last_updated_by, last_updated_at)
    VALUES ('Civil and Structure', 1,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Finishing Works', 1,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Plumbing and Firefighting', 1,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('HVAC', 1,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Electrical', 1,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Infrastructure', 0,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Landscaping', 1,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Renewable Energy', 0,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Facade works', 0,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Waterproofing', 0,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Construction Equipments', 1,0, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Safety', 1,0, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00'),
    ('Others', 0,1, 0, 'admin', '2020-01-01 10:00:00', 'admin', '2020-02-02 10:00:00');

	
	CREATE TABLE material_trade(
material_trade_id serial PRIMARY KEY,
trade_id integer REFERENCES trade(trade_id) NOT NULL,
material_id integer REFERENCES materials(id) NOT NULL,
status SMALLINT DEFAULT 0 NOT NULL,
created_by VARCHAR(50) NOT NULL,
created_at TIMESTAMP NOT NULL,
last_updated_by VARCHAR(50) NOT NULL,
last_updated_at TIMESTAMP NOT NULL
);
DROP INDEX material_code_material_group_idxs;
CREATE UNIQUE INDEX material_code_material_pid ON materials (material_code,pid);

alter table materials add column sort_seq int
update materials set material_unit='sqm' where material_unit='sq m';
update materials set material_unit='sqm' where material_unit='sq m';
update materials set material_unit='set' where material_unit='Set';
update materials set material_unit='nos' where material_unit='No.';
update materials set material_unit='LS' where material_unit='L.S.';
update materials set material_unit='litre' where material_unit='Litre';
update materials set material_unit='sqft' where material_unit='sq ft';
update materials set material_unit='nos' where material_unit='duplicate';
update materials set material_unit='rmt' where material_unit='R.m';
update materials set material_unit='packet' where material_unit='Packet';
update materials set material_unit='bag' where material_unit='Bag';
update materials set material_unit='roll' where material_unit='Roll';