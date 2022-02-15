/* use client address for receiver for dummy data */
USE `watch`;

DELETE FROM `payment_request`;

INSERT INTO `payment_request` VALUES
  (1,'btc',NULL,'tb1qh7etm6w40e66u4f05zks8m4kc5p76tdqa5fcle6lendush3sl7ts57d28e','client','tb1qnk5aanyf2nk0kx9he5k6n0akny2m9xkh0f6k6jnygevgf5rzpdsqk4hj4j',0.0001,false,now()),
  (2,'btc',NULL,'tb1qh7etm6w40e66u4f05zks8m4kc5p76tdqa5fcle6lendush3sl7ts57d28e','client','tb1qx5u64ydftwdlnqhrf9xujkjfq9ptqkngrpjkjdhr03f6j3hnwlxqpr6ngw',0.0002,false,now()),
  (3,'btc',NULL,'tb1qh7etm6w40e66u4f05zks8m4kc5p76tdqa5fcle6lendush3sl7ts57d28e','client','tb1q92g7wgfjmp4p0v9gsmx3ka3d2nr8e23unkekvwt6qduxt5say7tqszv6yj',0.0025,false,now()),
  (4,'btc',NULL,'tb1qh7etm6w40e66u4f05zks8m4kc5p76tdqa5fcle6lendush3sl7ts57d28e','client','tb1qfzkgrq6rqayvht7dmlps48pctrr0dpj3v0xhh9svkzt7l6mrqyts0yt3zh',0.0015,false,now()),
  (5,'btc',NULL,'tb1qh7etm6w40e66u4f05zks8m4kc5p76tdqa5fcle6lendush3sl7ts57d28e','client','tb1qhum482cj8wmh2j5yqu7y74fypgce37tdz0r5ksru83944p83g7nsq7macg',0.0022,false,now());
