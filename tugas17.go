package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Barang struct {
	Id          int
	NamaBarang  string
	HargaBarang int
}

func koneksi() (*sql.DB, error) {
	var username string = "root"
	var password string = ""
	var host string = "localhost"
	var database string = "balajar_go"
	db, err := sql.Open("mysql", fmt.Sprintf("%s@%stcp(%s:3306)/%s", username, password, host, database))
	if err != nil {
		return nil, err
	}
	return db, nil

}

func sql_tampil(w http.ResponseWriter, r *http.Request) {
	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	rows, err := db.Query("select * from tbl_barang")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	var result []Barang

	for rows.Next() {
		var each = Barang{}
		var err = rows.Scan(&each.Id, &each.NamaBarang, &each.HargaBarang)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	var results, _ = json.Marshal(result)

	w.Write(results)
}

func sql_cari(w http.ResponseWriter, r *http.Request) {
	var id = r.FormValue("id")

	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	rows, err := db.Query("select * from tbl_barang where id = ?", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	var result []Barang

	for rows.Next() {
		var each = Barang{}
		var err = rows.Scan(&each.Id, &each.NamaBarang, &each.HargaBarang)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	var results, _ = json.Marshal(result)

	w.Write(results)
}

func sql_tambah(w http.ResponseWriter, r *http.Request) {
	var nama = r.FormValue("nama_barang")
	var harga = r.FormValue("harga_barang")

	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	_, err = db.Exec("insert into tbl_barang value (?,?,?)", nil, nama, harga)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Fprintf(w, "%s ", "berhasil di tambahkan")

}

func sql_update(w http.ResponseWriter, r *http.Request) {
	var nama = r.FormValue("nama_barang")
	var harga = r.FormValue("harga_barang")
	var id = r.FormValue("id")

	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	_, err = db.Exec("UPDATE `tbl_barang` SET `nama_barang`= ?,`harga_barang`= ? WHERE id = ?", nama, harga, id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Fprintf(w, "%s ", "berhasil di update")

}

func sql_delete(w http.ResponseWriter, r *http.Request) {
	var id = r.FormValue("id")

	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM `tbl_barang` WHERE id = ?", id)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Fprintf(w, "%s ", "berhasil di hapus")

}

func main() {
	http.HandleFunc("/", sql_tampil)
	http.HandleFunc("/tambah", sql_tambah)
	http.HandleFunc("/update", sql_update)
	http.HandleFunc("/delete", sql_delete)
	http.HandleFunc("/cari", sql_cari)

	http.ListenAndServe(":8000", nil)
}
