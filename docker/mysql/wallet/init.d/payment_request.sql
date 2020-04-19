USE `wallet`;

/* use client address for receiver */

DELETE FROM `payment_request`;

INSERT INTO `payment_request` VALUES
  (1,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom1','2N62NC214arsygo3RN3vAYoJ1HL4aL9bn5d',0.001,false,now()),
  (2,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom2','2NAdK1eoi2spwggfScGxrkGDT9MFaSLZWez',0.002,false,now()),
  (3,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom3','2N1GqmQz9sSpAsS94AUjYb1AjMqDems6Z3p',0.0025,false,now()),
  (4,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom4','2NCGSUS2tGeGfut1hyySvYHtyBtYk3yX1rz',0.0015,false,now()),
  (5,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom5','2NAYAynzTtHDWqqQjzLcoSDTNtKAf7gakPw',0.0022,false,now());
