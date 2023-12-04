#![cfg_attr(not(feature = "std"), no_std, no_main)]

#[ink::contract]
mod hash_store {
    use ink::storage::Mapping;
    use ink_env::hash;
    use scale::{Decode, Encode};

    pub fn hash_keccak_256(input: &[u8]) -> [u8; 32] {
        let mut output = <hash::Keccak256 as hash::HashOutput>::Type::default();
        ink_env::hash_bytes::<hash::Keccak256>(input, &mut output);
        output
    }

    #[derive(Encode, Decode, Debug, Clone)]
    #[cfg_attr(
        feature = "std",
        derive(scale_info::TypeInfo, ink::storage::traits::StorageLayout)
    )]
    pub struct ChallengeKey {
        from_bank: [u8; 32], //Hash,
        to_bank: [u8; 32],   //Hash,
        block_number: BlockNumber,
    }

    impl ChallengeKey {
        fn new(from_bank: [u8; 32], to_bank: [u8; 32], block_number: BlockNumber) -> Self {
            Self {
                from_bank: from_bank,
                to_bank: to_bank,
                block_number: block_number,
            }
        }
    }

    #[derive(Encode, Decode, Debug, Clone)]
    #[cfg_attr(
        feature = "std",
        derive(scale_info::TypeInfo, ink::storage::traits::StorageLayout)
    )]
    pub struct Challenge {
        hash: [u8; 32],
        solved: bool,
    }

    impl Challenge {
        fn new(hashv: [u8; 32], solv: bool) -> Self {
            Self {
                hash: hashv,
                solved: solv,
            }
        }
    }

    #[ink(storage)]
    pub struct HashStorage {
        // Mapping to store hash based on block number
        hash_map: Mapping<ChallengeKey, Challenge>,
    }

    #[ink(event)]
    pub struct ChallengeStored {
        #[ink(topic)]
        challenge_key: ChallengeKey,
        hash: [u8; 32],
    }

    #[ink(event)]
    pub struct ChallengeSolved {
        #[ink(topic)]
        challenge_key: ChallengeKey,
        hash: [u8; 32],
        block_number: BlockNumber,
    }

    impl HashStorage {
        #[ink(constructor)]
        pub fn new() -> Self {
            let nhash_map = Mapping::new();
            Self {
                hash_map: nhash_map,
            }
        }

        #[ink(message)]
        pub fn storage_challenge(
            &mut self,
            from_bank: [u8; 32],
            to_bank: [u8; 32],
            hash: [u8; 32],
        ) {
            // Get the current block number
            let current_block_number = self.env().block_number();

            let key = ChallengeKey::new(from_bank, to_bank, current_block_number);
            let challenge = Challenge::new(hash, false);

            // Store the hash in the mapping
            self.hash_map.insert(key.clone(), &challenge);

            // Emit an event for the stored hash
            self.env().emit_event(ChallengeStored {
                challenge_key: key,
                hash,
            });
        }

        #[ink(message)]
        pub fn get_challenge(
            &self,
            from_bank: [u8; 32],
            to_bank: [u8; 32],
            block_number: BlockNumber,
        ) -> Challenge {
            // Retrieve the hash from the mapping
            let key = ChallengeKey::new(from_bank, to_bank, block_number);
            let c = self.hash_map.get(&key).expect("None found");
            c
        }

        #[ink(message)]
        pub fn solve_challenge(
            &mut self,
            from_bank: [u8; 32],
            to_bank: [u8; 32],
            block_number: BlockNumber,
            preimage: [u8; 32],
        ) -> bool {
            let key = ChallengeKey::new(from_bank, to_bank, block_number);
            let Challenge { hash: h, solved: s } = self.hash_map.get(&key).expect("None found");
            assert!(!s, "Already solved");
            let result = hash_keccak_256(&preimage);
            let solved = result == h;
            if solved {
                self.hash_map
                    .insert(key.clone(), &Challenge::new(h, solved));
                self.env().emit_event(ChallengeSolved {
                    challenge_key: key,
                    hash: h,
                    block_number: self.env().block_number(),
                });
            }

            solved
        }
    }
}
