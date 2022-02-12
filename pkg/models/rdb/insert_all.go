package models

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

// InsertAll inserts all rows with the specified column values, using an executor.
func (o AccountKeySlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(accountKeyColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			accountKeyAllColumns,
			accountKeyColumnsWithDefault,
			accountKeyColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `account_key` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(accountKeyType, accountKeyMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into account_key")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o AddressSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(addressColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			addressAllColumns,
			addressColumnsWithDefault,
			addressColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `address` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(addressType, addressMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into address")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o AuthAccountKeySlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(authAccountKeyColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			authAccountKeyAllColumns,
			authAccountKeyColumnsWithDefault,
			authAccountKeyColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `auth_account_key` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(authAccountKeyType, authAccountKeyMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into auth_account_key")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o AuthFullpubkeySlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(authFullpubkeyColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			authFullpubkeyAllColumns,
			authFullpubkeyColumnsWithDefault,
			authFullpubkeyColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `auth_fullpubkey` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into auth_fullpubkey")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o BTCTXSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {

		nzDefaults := queries.NonZeroDefaultSet(btcTXColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			btcTXAllColumns,
			btcTXColumnsWithDefault,
			btcTXColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `btc_tx` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(btcTXType, btcTXMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into btc_tx")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o BTCTXInputSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(btcTXInputColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			btcTXInputAllColumns,
			btcTXInputColumnsWithDefault,
			btcTXInputColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `btc_tx_input` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(btcTXInputType, btcTXInputMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into btc_tx_input")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o BTCTXOutputSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(btcTXOutputColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			btcTXOutputAllColumns,
			btcTXOutputColumnsWithDefault,
			btcTXOutputColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `btc_tx_output` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(btcTXOutputType, btcTXOutputMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into btc_tx_output")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o EthDetailTXSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {

		nzDefaults := queries.NonZeroDefaultSet(ethDetailTXColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			ethDetailTXAllColumns,
			ethDetailTXColumnsWithDefault,
			ethDetailTXColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `eth_detail_tx` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(ethDetailTXType, ethDetailTXMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into eth_detail_tx")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o PaymentRequestSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(paymentRequestColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			paymentRequestAllColumns,
			paymentRequestColumnsWithDefault,
			paymentRequestColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `payment_request` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(paymentRequestType, paymentRequestMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into payment_request")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o SeedSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(seedColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			seedAllColumns,
			seedColumnsWithDefault,
			seedColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `seed` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(seedType, seedMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into seed")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o TXSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(txColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			txAllColumns,
			txColumnsWithDefault,
			txColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `tx` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(txType, txMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into tx")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o XRPAccountKeySlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(xrpAccountKeyColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			xrpAccountKeyAllColumns,
			xrpAccountKeyColumnsWithDefault,
			xrpAccountKeyColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `xrp_account_key` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(xrpAccountKeyType, xrpAccountKeyMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into xrp_account_key")
	}

	return nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o XRPDetailTXSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {

		nzDefaults := queries.NonZeroDefaultSet(xrpDetailTXColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			xrpDetailTXAllColumns,
			xrpDetailTXColumnsWithDefault,
			xrpDetailTXColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `xrp_detail_tx` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(xrpDetailTXType, xrpDetailTXMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into xrp_detail_tx")
	}

	return nil
}
