注意点：

1，瑞波币货币金额【XRP的数量表示为字符串。（XRP的精度相当于64位整数，但JSON整数限制为32位，因此如果以JSON整数表示，XRP可能会溢出。）XRP在“drop”中正式指定，相当于0.000001（1每个XRP的百万分之一。因此，要在JSON文档中表示1.0 XRP，您可以编写："1000000"】

2，日期格式。博主是通过日期来保存用户交易记录处理进度的。有人用的是marker保存进度。

3，用户充值地址。博主通过在主地址上加用户id（address_userId）。告知用户充值地址。用户充值时需主动处理地址。在备注中写入userId。

package com.tn.web.service.coin;
import java.io.IOException;
import java.math.BigDecimal;
import java.text.MessageFormat;
import java.util.*;
 
 
import org.apache.commons.codec.binary.Base64;
import org.slf4j.Logger;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;
 
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.tn.base.Log;
import com.tn.constant.CoinConstant;
import com.tn.util.BigDecimalUtil;
import com.tn.util.DateUtil;
import com.tn.util.HttpUtil;
 
/**
 * XRP操作
 * @author xhl
 * @create 2017/10/27
 **/
@Service
public class CoinXrpService{
 
	//URL都为正式网络地址可信任
	private  String getUrl = "https://data.ripple.com";
	private  String postUrl = "https://s1.ripple.com:51234";
	private  String address = "";
	private  String password = "";
 
    private Logger log = Log.get();
 
    private final static String RESULT = "result";
    private final static String SUCCESS = "success";
    private final static String TES_SUCCESS = "tesSUCCESS";
    private final static String METHOD_GET_TRANSACTION = "/v2/accounts/{0}/transactions";
    private final static String METHOD_GET_BALANCE = "/v2/accounts/{0}/balances";
    private final static String METHOD_POST_SIGN = "sign";
    private final static String METHOD_POST_SUBMIT = "submit";
 
    public String getAddress(){
    	return address;
    }
    
    /**
     * 发送交易
     * @param address
     * @param value
     * @return
     */
    public String send(String toAddress,double value){
    	String txBlob = this.sign(toAddress, value);
    	if (StringUtils.isEmpty(txBlob)) {
			log.error("签名失败:{}",toAddress);
			return null;
		}
    	HashMap<String, Object> params = new HashMap<String, Object>();
    	params.put("tx_blob", txBlob);
    	//签名
    	JSONObject json = doRequest(METHOD_POST_SUBMIT, params);
    	if (!isError(json)) {
    		JSONObject result = json.getJSONObject(RESULT);
    		if (result != null) {
				if (TES_SUCCESS.equals(result.getString("engine_result"))) {
					String hash = result.getJSONObject("tx_json").getString("hash");
					if (!StringUtils.isEmpty(hash)) {
						log.info("转账成功：toAddress:{},value:{},hash:{}",toAddress,value,hash);
						return hash;
					}else {
						log.error("转账失败：toAddress:{},value:{},hash:{}",toAddress,value,hash);
					}
				}
			}
		}
    	return null;
    }
 
    /**
     * 签名
     * @param address
     * @param value
     * @return tx_blob
     */
    public String sign(String toAddress,Double value){
    	//瑞波币余额存储加六位长度
    	value = BigDecimalUtil.mul(value, 1000000);
    	Integer vInteger = BigDecimal.valueOf(value).intValue();
    	JSONObject txJson = new JSONObject();
    	txJson.put("Account", address);
    	txJson.put("Amount", vInteger.toString());
    	txJson.put("Destination", toAddress);
    	txJson.put("TransactionType", "Payment");
    	txJson.put("DestinationTag", "1");
    	HashMap<String, Object> params = new HashMap<String, Object>();
    	params.put("secret", password);
    	params.put("tx_json", txJson);
    	params.put("offline", false);
    	//签名
    	JSONObject json = doRequest(METHOD_POST_SIGN, params);
    	if (!isError(json)) {
    		JSONObject result = json.getJSONObject(RESULT);
    		if (result != null) {
				if (SUCCESS.equals(result.getString("status"))) {
					return result.getString("tx_blob");
				}
			}
		}
    	return null;
    }
    
    /**
     * XRP查询余额
     * @return
     */
    public double getBalance(){
    	HashMap<String, String> params = new HashMap<String, String>();
    	params.put("currency", CoinConstant.COIN_XRP);
    	String re = HttpUtil.jsonGet(getUrl + MessageFormat.format(METHOD_GET_BALANCE, address), params);
        log.info("获取XRP余额:{}",re);
        if(!StringUtils.isEmpty(re)){
            JSONObject json = JSON.parseObject(re);
            if (SUCCESS.equals(json.getString(RESULT))) {
            	JSONArray array = json.getJSONArray("balances");
            	if (array != null && array.size() > 0) {
            		//总余额
            		double balance = array.getJSONObject(0).getDoubleValue("value");
            		if (balance >= 20) {
            			//可用余额    xrp会冻结20个币
						return BigDecimalUtil.sub(balance, 20);
					}
				}
    		}
        }
        return 0.00;
    }
    
    public Long parseTransaction(String startTm) {
    	HashMap<String, String> params = new HashMap<String, String>();
    	if (!StringUtils.isEmpty(startTm)) {
    		//记录时间格式目前XRP精确到秒
    		Date d = new Date(BigDecimalUtil.longAdd(Long.parseLong(startTm), 1000L));
        	params.put("start", DateUtil.dateToString(d, "yyyy-MM-dd'T'HH:mm:ss"));
		}
    	params.put("result", "tesSUCCESS");
    	params.put("type", "Payment");
    	params.put("limit", "100");
    	String re = HttpUtil.jsonGet(getUrl + MessageFormat.format(METHOD_GET_TRANSACTION, address), params);
        if(!StringUtils.isEmpty(re)){
            JSONObject json = JSON.parseObject(re);
            if (SUCCESS.equals(json.getString(RESULT))) {
//            	marker = json.getString("marker");
            	JSONArray transactions = json.getJSONArray("transactions");
            	if (transactions != null && transactions.size() > 0) {
					for (Object object : transactions) {
						JSONObject transaction = (JSONObject)object;
						String hash = transaction.getString("hash");
						String dateString = transaction.getString("date");
						Date date =DateUtil.getStringToDate(dateString,"yyyy-MM-dd'T'HH:mm:ss");
						JSONObject tx = transaction.getJSONObject("tx");
						String destinationTag = tx.getString("DestinationTag");
						if (StringUtils.isEmpty(destinationTag)) {
							log.info("非用户充值记录");
		                	return date.getTime();
						}
						String to = tx.getString("Destination");
						if (!address.equals(to)) {
							log.info("非用户充值记录,地址不一致");
		                	return date.getTime();
						}
						//根据tag查找用户ID
						/*UserEntity user = userService.getUserById(Integer.parseInt(destinationTag));
						if (user == null) {
							log.info("用户不存在：{}",destinationTag);
		                	return date.getTime();
						}*/
						double amount = tx.getDoubleValue("Amount");
						if (amount > 0 ) {
							amount = BigDecimalUtil.div(amount, 1000000, 6);
						}else {
							log.error("交易金额异常：{}",amount);
		                	return date.getTime();
						}
						//UserCoinRecordEntity record = new UserCoinRecordEntity();
		                /*record.setCoinType(CoinConstant.COIN_XRP);
		                record.setAddress(to+"_"+destinationTag);
		                record.setTxid(hash);
		                record.setUserId(user.getUserId());*/
		                try {
		                    //rechargeParse(record);
		                	return date.getTime();
		                }catch (DuplicateKeyException e){
		                    log.error("eos hash:{} userid:{}  coin:{} 这个区块已经处理了",hash,"XRP");
		                    return null;
		                }
					}
				}
    		}
        }
        return null;
    }
 
    private boolean isError(JSONObject json){
        if( json == null || (!StringUtils.isEmpty(json.getString("error")) && json.get("error") != "null")){
            return true;
        }
        return false;
    }
    private JSONObject doRequest(String method,Object... params){
        JSONObject param = new JSONObject();
        param.put("id",System.currentTimeMillis()+"");
        param.put("jsonrpc","2.0");
        param.put("method",method);
        if(params != null){
            param.put("params",params);
        }
        String creb = Base64.encodeBase64String((address+":"+password).getBytes());
        Map<String,String> headers = new HashMap<>(2);
        headers.put("Authorization","Basic "+creb);
        String resp = "";
        try{
            resp = HttpUtil.jsonPost(postUrl,headers,param.toJSONString());
        }catch (Exception e){
        	log.info(e.getMessage());
            if (e instanceof IOException){
                resp = "{}";
            }
        }
        log.info(resp);
        return JSON.parseObject(resp);
    }
 
}
处理用户充值：定时任务扫描

    /**
     * XRP处理
     */
    private void xrpJob(){
        //获取数据库保存的记录处理进度（日期）
        String startTm = coinParseService.getBlockHeight(CoinConstant.COIN_XRP);
        log.info("XRP当前处理进度：{}",startTm);
        Long start = coinXrpService.parseTransaction(startTm);
        if (!StringUtils.isEmpty(start)) {
            log.info("xrp执行完毕");
            coinParseService.updateBlockRecord(CoinConstant.COIN_XRP, start.toString());
		}
    }
希望能帮到大家，欢迎大家一起分享。

觉得有用请打赏，你的鼓励就是我的动力！
--------------------- 
作者：其修远兮 
来源：CSDN 
原文：https://blog.csdn.net/liu1765686161/article/details/82492937 
版权声明：本文为博主原创文章，转载请附上博文链接！