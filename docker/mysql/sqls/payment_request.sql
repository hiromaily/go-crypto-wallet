/* use client address for receiver */

DELETE FROM `payment_request`;

INSERT INTO `payment_request` VALUES
  (1,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom1','2NBqM5FFAW9N9yus1rEVKoMAujCEexbGzpf',0.001,false,now()),
  (2,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom2','2N7whxoFvzHjMFMdxfBuprzcjt2ijgrrpyT',0.002,false,now()),
  (3,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom3','2NCZEK64JomPMsabwibg4ZaYdctjjJuEXvd',0.0025,false,now()),
  (4,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom4','2MycJFujPwBWrNjtFT28XXZAKhypXrSwRTk',0.0015,false,now()),
  (5,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','tom5','2N1rRRZ3xFxcB17y8Bbg9AtMawkiGsTwK7D',0.0022,false,now());
