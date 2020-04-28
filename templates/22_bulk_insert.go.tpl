{{- $alias := .Aliases.Table .Table.Name -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o {{$alias.UpSingular}}Slice) InsertAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, columns boil.Columns) error {
    ln := int64(len(o))
    if ln == 0 {
        return nil
    }

    var sql string
    vals := []interface{}{}
    for i, row := range o {
        {{- template "timestamp_bulk_insert_helper" . }}

        {{if not .NoHooks -}}
        if err := row.doBeforeInsertHooks(ctx, exec); err != nil {
            return err
        }
        {{- end}}

        nzDefaults := queries.NonZeroDefaultSet({{$alias.DownSingular}}ColumnsWithDefault, row)
        wl, _ := columns.InsertColumnSet(
            {{$alias.DownSingular}}AllColumns,
            {{$alias.DownSingular}}ColumnsWithDefault,
            {{$alias.DownSingular}}ColumnsWithoutDefault,
            nzDefaults,
        )
        if i == 0 {
            sql = "INSERT INTO {{$schemaTable}} " + "({{.LQ}}" + strings.Join(wl, "{{.RQ}},{{.LQ}}") + "{{.RQ}})" + " VALUES "
        }
        sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
        if i != len(o)-1 {
            sql += ","
        }
        valMapping, err := queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, wl)
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

    {{if .NoContext -}}
    _, err := exec.Exec(ctx, sql, vals...)
    {{else -}}
    _, err := exec.ExecContext(ctx, sql, vals...)
    {{end -}}

    if err != nil {
        return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
    }

    return nil
}

{{- define "timestamp_bulk_insert_helper" -}}
    {{- if not .NoAutoTimestamps -}}
    {{- $colNames := .Table.Columns | columnNames -}}
    {{if containsAny $colNames "created_at" "updated_at"}}
        {{if not .NoContext -}}
    if !boil.TimestampsAreSkipped(ctx) {
        {{end -}}
        currTime := time.Now().In(boil.GetLocation())
        {{range $ind, $col := .Table.Columns}}
            {{- if eq $col.Name "created_at" -}}
                {{- if eq $col.Type "time.Time" }}
        if row.CreatedAt.IsZero() {
            row.CreatedAt = currTime
        }
                {{- else}}
        if queries.MustTime(row.CreatedAt).IsZero() {
            queries.SetScanner(&row.CreatedAt, currTime)
        }
                {{- end -}}
            {{- end -}}
            {{- if eq $col.Name "updated_at" -}}
                {{- if eq $col.Type "time.Time"}}
        if row.UpdatedAt.IsZero() {
            row.UpdatedAt = currTime
        }
                {{- else}}
        if queries.MustTime(row.UpdatedAt).IsZero() {
            queries.SetScanner(&row.UpdatedAt, currTime)
        }
                {{- end -}}
            {{- end -}}
        {{end}}
        {{if not .NoContext -}}
    }
        {{end -}}
    {{end}}
    {{- end}}
{{- end -}}
