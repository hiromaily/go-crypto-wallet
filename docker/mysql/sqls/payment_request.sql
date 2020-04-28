/* use client address for receiver for dummy data */

DELETE FROM `payment_request`;

INSERT INTO `payment_request` VALUES
  (1,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','client','2NBqM5FFAW9N9yus1rEVKoMAujCEexbGzpf',0.001,false,now()),
  (2,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','client','2N7whxoFvzHjMFMdxfBuprzcjt2ijgrrpyT',0.002,false,now()),
  (3,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','client','2NCZEK64JomPMsabwibg4ZaYdctjjJuEXvd',0.0025,false,now()),
  (4,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','client','2MycJFujPwBWrNjtFT28XXZAKhypXrSwRTk',0.0015,false,now()),
  (5,'btc',NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','client','2N1rRRZ3xFxcB17y8Bbg9AtMawkiGsTwK7D',0.0022,false,now());
