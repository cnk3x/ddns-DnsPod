package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Status struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	CreateAt string `json:"created_at"`
}

type Domain struct {
	Id         int    `json:"id"`
	DomainName string `json:"punyCode"`
}

type Record struct {
	RecordId    string `json:"id"`
	RecordName  string `json:"name"`
	RecordType  string `json:"type"`
	RecordValue string `json:"value"`
	LoginKey    string
	DomainId    string
	DomainName  string
}

type DomainResult struct {
	Status  Status   `json:"status"`
	Domains []Domain `json:"domains"`
}
type RecordResult struct {
	Status  Status   `json:"status"`
	Records []Record `json:"records"`
}
type ModifyResult struct {
	Status Status `json:"status"`
}

func GetIp() (string, error) {
	response, err := http.Get("http://ip.cip.cc")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func GetDpd(domainName string, token string) (*Domain, error) {
	v := url.Values{}
	v.Set("offset", "0")
	v.Set("length", "1")
	v.Set("keyword", domainName)
	body, err := dpPost("https://dnsapi.cn/Domain.List", v, token)
	if err != nil {
		return nil, err
	}

	Debug("GetDpd:" + token + "," + domainName + "," + string(body))

	var domainResult DomainResult
	err = json.Unmarshal(body, &domainResult)
	if err != nil {
		return nil, err
	}
	if len(domainResult.Domains) > 0 {
		dpd := domainResult.Domains[0]
		return &dpd, nil
	}
	return nil, errors.New("没有记录:" + domainName)
}

func GetDpr(domainId string, name string, token string) (*Record, error) {
	v := url.Values{}
	v.Set("domain_id", domainId)
	v.Set("length", "10")
	v.Set("keyword", name)
	body, err := dpPost("https://dnsapi.cn/Record.List", v, token)
	if err != nil {
		return nil, err
	}

	Debug("GetDpr:" + domainId + "," + name + "," + token + "\n" + string(body))

	var record RecordResult
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil, err
	}
	if len(record.Records) > 0 {
		for _, r := range record.Records {
			if r.RecordType == "A" && strings.ToLower(r.RecordName) == strings.ToLower(name) {
				return &r, nil
			}
		}
	}
	return nil, errors.New("没有记录:" + domainId + "," + name)
}

func DoUpdate(domainId string, recordId string, subDomain string, ip string, token string) (string, error) {
	v := url.Values{}
	v.Set("domain_id", domainId)
	v.Set("record_id", recordId)
	v.Set("sub_domain", subDomain)
	v.Set("record_line", "默认")
	v.Set("value", ip)

	Debug("domainId:" + domainId)
	Debug("recordId:" + recordId)

	body, err := dpPost("https://dnsapi.cn/Record.Ddns", v, token)
	if err != nil {
		return "", err
	}
	Debug("DoUpdate:" + string(body))

	var m_result ModifyResult
	err = json.Unmarshal(body, &m_result)
	if err != nil {
		return "", err
	}

	if m_result.Status.Code == "1" {
		return m_result.Status.Message, nil
	} else {
		return m_result.Status.Message, errors.New(m_result.Status.Message)
	}
}

func FixAdd(records []*Record, domainName string, recordName string, token string) ([]*Record, error) {
	record := &Record{DomainName: domainName, RecordName: recordName, LoginKey: token}
	d, d_err := GetDpd(record.DomainName, record.LoginKey)
	if d_err == nil {
		record.DomainId = strconv.Itoa(d.Id)
		r, r_err := GetDpr(record.DomainId, record.RecordName, record.LoginKey)
		if r_err == nil {
			record.LoginKey = record.LoginKey
			record.RecordId = r.RecordId
			record.RecordName = r.RecordName
			record.RecordType = r.RecordType
			record.RecordValue = r.RecordValue
		} else {
			return records, errors.New("R:" + r_err.Error())
		}
	} else {
		return records, errors.New("D:" + d_err.Error())
	}
	return append(records, record), nil
}

func dpPost(addr string, v url.Values, loginKey string) ([]byte, error) {
	var body []byte
	v.Set("login_token", loginKey)
	v.Set("format", "json")
	v.Set("lang", "cn")
	v.Set("error_on_empty", "no")
	reqest, err := http.NewRequest("POST", addr, strings.NewReader(v.Encode()))
	if err != nil {
		return body, err
	}

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Add("User-Agent", "Hao DDns/1.0.0 (yao@wenaiyao.com)")
	reqest.Header.Add("Accept", "*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Add("Connection", "keep-alive")

	client := &http.Client{}

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		return body, err
	}

	if response.StatusCode == 200 {
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			var reader *gzip.Reader
			reader, _ = gzip.NewReader(response.Body)
			body, _ = ioutil.ReadAll(reader)
		default:
			body, _ = ioutil.ReadAll(response.Body)
		}
		return body, nil
	} else {
		return body, errors.New(strconv.Itoa(response.StatusCode) + ":" + response.Status)
	}

	return body, errors.New("未知错误.")
}
