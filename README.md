# delivery
микросервис доставки

## Запросы к БД
### Выборки
```sql
SELECT * FROM public.couriers;
SELECT * FROM public.transports;
SELECT * FROM public.orders;

SELECT * FROM public.outbox;
```

### Очистка БД (все кроме справочников)
```sql
DELETE FROM public.couriers;
DELETE FROM public.transports;
DELETE FROM public.orders;
DELETE FROM public.outbox;
```


### Добавить курьеров

#### Пеший
```sql
INSERT INTO public.transports(
id, name, speed)
VALUES ('921e3d64-7c68-45ed-88fb-97ceb8148a7e', 'Пешком', 1);
INSERT INTO public.couriers(
id, name, transport_id, location_x, location_y, status)
VALUES ('bf79a004-56d7-4e5f-a21c-0a9e5e08d10d', 'Пеший', '921e3d64-7c68-45ed-88fb-97ceb8148a7e', 1, 3, 'free');
```

#### Вело
```sql
INSERT INTO public.transports(
id, name, speed)
VALUES ('b96a9d83-aefa-4d06-99fb-e630d17c3868', 'Велосипед', 2);
INSERT INTO public.couriers(
id, name, transport_id, location_x, location_y, status)
VALUES ('db18375d-59a7-49d1-bd96-a1738adcee93', 'Вело', 'b96a9d83-aefa-4d06-99fb-e630d17c3868', 4,5, 'free');
```

#### Авто
```sql
INSERT INTO public.transports(
id, name, speed)
VALUES ('c24d3116-a75c-4a4b-9b22-1a7dc95a8c79', 'Машина', 3);
INSERT INTO public.couriers(
id, name, transport_id, location_x, location_y, status)
VALUES ('407f68be-5adf-4e72-81bc-b1d8e9574cf8', 'Авто', 'c24d3116-a75c-4a4b-9b22-1a7dc95a8c79', 7,9, 'free');     
```