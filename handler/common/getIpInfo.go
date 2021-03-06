package common

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"project-name/handler"
	"project-name/handler/common/dto"
	"project-name/util/byteDecoder"
)

const baseUrl = "http://whois.pconline.com.cn/ipJson.jsp?ip=%s&json=true&level=2"

// @Summary 获取ip归属地信息
// @Description desc
// @Tags common
// @Accept json
// @Produce json
// @Param ip path string true "desc" #where_in:path/query/body
// @Success 200 {object} dto.GetIpInfoDto ""
// @Router /v1/common/ip/:ip [get] #method:get/post/put/delete
func GetIpInfo(c *gin.Context) {
	var (
		ip        string
		ipInfo    *dto.GetIpInfoDto
		err       error
		retryTime int = 3
	)

	ip = c.Param("ip")
	for i := 0; i < retryTime; i++ {
		ipInfo, err = getIpInfoReal(ip)
		if err == nil {
			break
		}
	}

	if err != nil {
		handler.Response(c, err, nil)
		return
	}

	handler.Response(c, nil, ipInfo)
}

func getIpInfoReal(ip string) (ipInfo *dto.GetIpInfoDto, err error) {
	var (
		url      string
		resp     *http.Response
		respBody []byte
	)

	url = fmt.Sprintf(baseUrl, ip)
	ipInfo = &dto.GetIpInfoDto{Ip: ip,}

	// get ip info
	resp, err = http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 编码转换
	respBody, err = byteDecoder.Decode(respBody, byteDecoder.GBK)
	if err != nil {
		panic(err)
	}
	
	err = json.Unmarshal(respBody, &ipInfo)
	if err != nil {
		panic(err)
	}

	return
}
