USE `keygen`;

INSERT INTO `account_key` (`coin`, `account`, `p2pkh_address`, `p2sh_segwit_address`, `bech32_address`, `full_public_key`, `multisig_address`, `redeem_script`, `wallet_import_format`, `idx`, `addr_status`, `updated_at`) VALUES
  ('eth', 'payment', '0xF7763dFB4eDeCd5854125c8dCa3531aed8077e55', '', '', '', '', '', '0xbd37a7868fe63209d7a6d4eea5991ea222df584481fb63270848ad6f02053ffc', 0, 0, now()),
  ('eth', 'deposit', '0x207749b4Ff22aB3B2c6D7c66414977B05957eDeB', '', '', '', '', '', '0xd0c6093ed9ee7e242fb52755412b1f1ece96b3d27ce7b6d2452262e6b2ae8b0e', 0, 0, now()),
  ('eth', 'stored', '0xa804cb16847Fb1D74677a726530F7C48cffA194E', '', '', '', '', '', '0xa5ad53b4d9a4ff2e385d6ee607b013f5bdb849d590255ffc347a7411a4d03493', 0, 0, now()),
  ('eth', 'client', '0xEd3FE4d0b757916AA435bAf6340608A98A508289', '', '', '', '', '', '0xce0b404c072a0920efc02302f674c66a0693e2d8ffa9e5cd06b71c759efc648c', 0, 0, now()),
  ('eth', 'client', '0xAB0009d07741319623681f30654369d8ce0F0961', '', '', '', '', '', '0x1a476014eb5b2bbc8736a7956eb72669a67f5b35c92683ab5057f2bae378232f', 0, 0, now()),
  ('eth', 'client', '0x08AC81d42a4CeF69616be8cD0A2c2A5F2B3AC9af', '', '', '', '', '', '0x2f41eb7705401603f3e2b338e6220494a3ea0a389f61e0c5a3989b8e2bfc5b57', 0, 0, now()),
  ('eth', 'client', '0x39333038d15E37473540CF9cf9E128ba6B4B7C64', '', '', '', '', '', '0x6a79686ca2d79c57091a57a59abdc76d55f90b077e978fe9b1fc278437f6df73', 0, 0, now()),
  ('eth', 'client', '0xF4dC65E0C7aF201f7715b080b0259D8456B72469', '', '', '', '', '', '0x6ab9bc99c1b2f1a49ff731b4c3026a6b91404d25ccb814e908bf79b85f5eb6ca', 0, 0, now()),
  ('eth', 'client', '0x315065d24B5287ec85B62c72eda525d57f2a661e', '', '', '', '', '', '0x7b077584258e2cb7fc87e90c057cc028bd2f782163c4291c41b163c55e8eccad', 0, 0, now()),
  ('eth', 'client', '0xA27c459AC3F644b7fBb0d2DCb29d7E9EdDe1c4f1', '', '', '', '', '', '0x90f6d5f32904eeba162063b2645b30f060293cdd9f6a4d826328358a7e5f16fd', 0, 0, now());
