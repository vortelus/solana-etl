package main

import (
	"fmt"
	"strconv"
)

func transformToAccountRecords(rawBlock *EtlBlock) ([]*AccountRecord, []*TokenRecord) {
	var accountRecords []*AccountRecord
	var tokenRecords []*TokenRecord

	blockSlot := int64(rawBlock.Slot)
	blockHash := rawBlock.TableContext.BlockHash

	var blockTimestamp *int64
	if rawBlock.TableContext.BlockTimestamp != nil {
		_timestamp := rawBlock.TableContext.BlockTimestamp.Timestamp * 1e6
		blockTimestamp = &_timestamp
	}
	// first accounts, and later tokens
	for _, accounts_data := range rawBlock.Accounts {
		for _, account := range accounts_data.GetAccounts() {
			var retrievalTimestamp *int64
			if account.RetrievalTimestamp != nil {
				_timestamp := account.RetrievalTimestamp.Timestamp * 1e6
				retrievalTimestamp = &_timestamp
			}

			authorizedVoters := make([]*AuthorizedVoterRecord, 0, len(account.AuthorizedVoters))
			for _, authorizedVoter := range account.AuthorizedVoters {
				if authorizedVoter != nil {
					authorizedVoters = append(authorizedVoters, &AuthorizedVoterRecord{
						AuthorizedVoter: &authorizedVoter.AuthorizedVoter,
						Epoch:           &authorizedVoter.Epoch,
					})
				}
			}

			if len(authorizedVoters) == 0 {
				authorizedVoters = append(authorizedVoters, &AuthorizedVoterRecord{AuthorizedVoter: nil, Epoch: nil})
			}

			priorVoters := make([]*PriorVoterRecord, 0, len(account.PriorVoters))
			for _, priorVoter := range account.PriorVoters {
				if priorVoter != nil {
					epochOfLastAuthorizedSwitch := int64(priorVoter.EpochOfLastAuthorizedSwitch)
					targetEpoch := int64(priorVoter.TargetEpoch)
					priorVoters = append(priorVoters, &PriorVoterRecord{
						AuthorizedPubkey:            &priorVoter.AuthorizedPubkey,
						EpochOfLastAuthorizedSwitch: &epochOfLastAuthorizedSwitch,
						TargetEpoch:                 &targetEpoch,
					})
				}
			}

			if len(priorVoters) == 0 {
				priorVoters = append(priorVoters, &PriorVoterRecord{AuthorizedPubkey: nil, EpochOfLastAuthorizedSwitch: nil, TargetEpoch: nil})
			}

			epochCredits := make([]*EpochCreditRecord, 0, len(account.EpochCredits))
			if account.Votes != nil {
				for _, epochCredit := range account.EpochCredits {
					if epochCredit != nil {
						epochCredits = append(epochCredits, &EpochCreditRecord{Credits: &epochCredit.Credits, Epoch: &epochCredit.Epoch, PreviousCredits: &epochCredit.PreviousCredits})
					}
				}
			}
			if len(epochCredits) == 0 {
				epochCredits = append(epochCredits, &EpochCreditRecord{Credits: nil, Epoch: nil, PreviousCredits: nil})
			}

			votes := make([]*VoteRecord, 0, len(account.Votes))
			if account.Votes != nil {
				for _, vote := range account.Votes {
					if vote != nil {
						votes = append(votes, &VoteRecord{
							ConfirmationCount: &vote.ConfirmationCount,
							Slot:              &vote.Slot,
						})
					}
				}
			}

			if len(votes) == 0 {
				votes = append(votes, &VoteRecord{ConfirmationCount: nil, Slot: nil})
			}

			var lastTimestamp []*TimestampRecord
			if account.LastTimestamp == nil {
				lastTimestamp = append(lastTimestamp, &TimestampRecord{Timestamp: nil, Slot: nil})
			} else {
				_timestamp := account.LastTimestamp.Timestamp * 1e6
				lastTimestamp = append(lastTimestamp, &TimestampRecord{Timestamp: &_timestamp, Slot: &account.LastTimestamp.Slot})
			}

			var dataRecord []*DataRecord
			if account.Data == nil {
				dataRecord = append(dataRecord, &DataRecord{Raw: nil, Encoding: nil})
			} else {
				dataRecord = append(dataRecord, &DataRecord{Raw: &account.Data.Raw, Encoding: &account.Data.Encoding})
			}

			rentEpoch := int64(account.RentEpoch)
			var tokenAmount *uint64
			if account.TokenAmount != nil && *account.TokenAmount != "" {
				_tokenAmount, err := strconv.ParseUint(*account.TokenAmount, 10, 64)
				if err != nil {
					fmt.Println("Error:", err)
					panic("terminating...")
				} else {
					tokenAmount = &_tokenAmount
				}
			}
			account_record := AccountRecord{
				BlockSlot:            &blockSlot,
				BlockTimestamp:       blockTimestamp,
				BlockHash:            &blockHash,
				TxSignature:          accounts_data.TxSignature,
				RetrievalTimestamp:   retrievalTimestamp,
				Pubkey:               &account.Pubkey,
				Executable:           &account.Executable,
				Lamports:             &account.Lamports,
				Owner:                account.Owner,
				RentEpoch:            &rentEpoch,
				Program:              account.Program,
				Space:                account.Space,
				AccountType:          account.AccountType,
				IsNative:             account.IsNative,
				Mint:                 account.Mint,
				State:                account.State,
				TokenAmount:          tokenAmount,
				TokenAmountDecimals:  account.TokenAmountDecimals,
				ProgramData:          account.ProgramData,
				AuthorizedVoters:     authorizedVoters,
				AuthorizedWithdrawer: account.AuthorizedWithdrawer,
				PriorVoters:          priorVoters,
				NodePubkey:           account.NodePubkey,
				Commission:           account.Commission,
				EpochCredits:         epochCredits,
				Votes:                votes,
				RootSlot:             account.RootSlot,
				LastTimestamp:        lastTimestamp,
				Data:                 dataRecord,
			}
			accountRecords = append(accountRecords, &account_record)
		}

		for _, token := range accounts_data.GetTokens() {
			var retrievalTimestamp *int64
			if token.RetrievalTimestamp != nil {
				_timestamp := token.RetrievalTimestamp.Timestamp * 1e6
				retrievalTimestamp = &_timestamp
			}

			var creators []*CreatorRecord
			if token.Creators == nil {
				creators = append(creators, &CreatorRecord{Address: nil, Verified: nil, Share: nil})
			} else {
				for _, creator := range token.Creators {
					share := int64(creator.Share)
					creators = append(creators, &CreatorRecord{Address: &creator.Address, Verified: &creator.Verified, Share: &share})
				}
			}

			tokenRecord := TokenRecord{
				BlockSlot:            &blockSlot,
				BlockTimestamp:       blockTimestamp,
				BlockHash:            &blockHash,
				TxSignature:          accounts_data.TxSignature,
				RetrievalTimestamp:   retrievalTimestamp,
				IsNft:                &token.IsNft,
				Mint:                 &token.Mint,
				UpdateAuthority:      &token.UpdateAuthority,
				Name:                 &token.Name,
				Symbol:               &token.Symbol,
				Uri:                  &token.Uri,
				SellerFeeBasisPoints: &token.SellerFeeBasisPoints,
				Creators:             creators,
				PrimarySaleHappened:  &token.PrimarySaleHappened,
				IsMutable:            &token.IsMutable,
			}
			tokenRecords = append(tokenRecords, &tokenRecord)

		}
	}
	return accountRecords, tokenRecords
}

func int64ToUint64(ptr *int64) *uint64 {
	u := uint64(*ptr)
	return &u
}
