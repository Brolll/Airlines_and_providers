CREATE TABLE "Provider" (
                            "id" varchar(2) PRIMARY KEY,
                            "name" varchar,
                            "code" SERIAL UNIQUE
);

CREATE TABLE "SSchema" (
                          "id" SERIAL PRIMARY KEY,
                          "name" varchar
);

CREATE TABLE "Airline" (
                           "id" varchar(2) PRIMARY KEY not null ,
                           "name" varchar not null

);

CREATE TABLE "Account" (
                           "id" SERIAL PRIMARY KEY,
                           "SchemaId" INT,
                           "name" varchar,
                           FOREIGN KEY ("SchemaId") REFERENCES "SSchema"(id)
);

create table "ProviderSchema" (
                            "schema_id" int,
                            "provider_id" varchar(2),
                            foreign key (schema_id)   references "SSchema"(id),
                            foreign key (provider_id) references "Provider"(id),
                            primary key (schema_id, provider_id)
);

CREATE table "ProviderAirline"(
                            "provider_code" int,
                            "airline_id" varchar(2),
                            foreign key (airline_id)   references "Airline"(id),
                            foreign key (provider_code)  references "Provider"(code),
                            primary key (airline_id, provider_code)
);



INSERT INTO "Airline"(id, name) VALUES ('SU', 'Аэрофлот');
INSERT INTO "Airline"(id, name) VALUES ('S7', 'S7');
INSERT INTO "Airline"(id, name) VALUES ('KV', 'КрасАвиа');
INSERT INTO "Airline"(id, name) VALUES ('U6', 'Уральские авиалинии');
INSERT INTO "Airline"(id, name) VALUES ('UT', 'ЮТэйр');
INSERT INTO "Airline"(id, name) VALUES ('FZ', 'Flydubai');
INSERT INTO "Airline"(id, name) VALUES ('JB', 'JetBlue');
INSERT INTO "Airline"(id, name) VALUES ('SJ', 'SuperJet');
INSERT INTO "Airline"(id, name) VALUES ('WZ', 'Wizz Air');
INSERT INTO "Airline"(id, name) VALUES ('N4', 'Nordwind Airlines');
INSERT INTO "Airline"(id, name) VALUES ('5N', 'SmartAvia');



INSERT INTO "Provider"(id, name) VALUES ('AA', 'AmericanAir');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (1,'FZ');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (1,'JB');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (1,'SJ');

INSERT INTO "Provider"(id, name) VALUES ('IF', 'InternationFlights');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (2,'SU');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (2,'S7');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (2,'FZ');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (2,'N4');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (2,'JB');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (2,'WZ');


INSERT INTO "Provider"(id, name) VALUES ('RS', 'RedStar');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (3,'SU');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (3,'S7');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (3,'KV');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (3,'U6');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (3,'UT');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (3,'N4');
INSERT INTO "ProviderAirline"(provider_code, airline_id) VALUES (3,'5N');



INSERT INTO "SSchema" (name) VALUES ('Main');
INSERT INTO "ProviderSchema" (schema_id, provider_id) VALUES (1, 'AA');
INSERT INTO "ProviderSchema" (schema_id, provider_id) VALUES (1, 'IF');
INSERT INTO "ProviderSchema" (schema_id, provider_id) VALUES (1, 'RS');

INSERT INTO "SSchema" (name) VALUES ('Test');
INSERT INTO "ProviderSchema" (schema_id, provider_id) VALUES (2, 'IF');
INSERT INTO "ProviderSchema" (schema_id, provider_id) VALUES (2, 'RS');

INSERT INTO "Account"("name","SchemaId") values ('Demo', 2);
INSERT INTO "Account"("name","SchemaId") values ('Develop', 2);
INSERT INTO "Account"("name","SchemaId") values ('Main', 1);



