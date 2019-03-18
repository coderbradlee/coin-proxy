具体注意点请参考:https://blog.csdn.net/liu1765686161/article/details/82492937

1.需要安装两个jar包   https://github.com/ripple-unmaintained/ripple-lib-java  中

ripple-bouncycastle 和  ripple-core .

2. 可以编译后放入maven仓库,然后引用 

例如:

<dependency>
            <groupId>com.ripple</groupId>
            <artifactId>ripple-bouncycastle</artifactId>
            <version>1.0.0</version>
        </dependency>
        <dependency>
          <groupId>com.ripple</groupId>
          <artifactId>ripple-core</artifactId>
          <version>1.0.0</version>
        </dependency>
3,直接上代码 

package com.tn.set.service.coin;
 
import java.io.IOException;
import java.math.BigDecimal;
import java.text.MessageFormat;
import java.util.*;
 
 
import org.apache.commons.codec.binary.Base64;
import org.slf4j.Logger;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;
 
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.ripple.core.coretypes.AccountID;
import com.ripple.core.coretypes.Amount;
import com.ripple.core.coretypes.uint.UInt32;
import com.ripple.core.types.known.tx.signed.SignedTransaction;
import com.ripple.core.types.known.tx.txns.Payment;
import com.tn.base.Log;
import com.tn.util.BigDecimalUtil;
import com.tn.util.DateUtil;
import com.tn.util.HttpUtil;
 
/**
 * XRP操作
 * @author xhl
 * @create 2017/10/27
 **/
@Service
public class CopyOfCoinXrpService {
 
	private String getUrl = "https://data.ripple.com";
	private String postUrl = "https://s1.ripple.com:51234";
	private String address = "rani9PZFVtQtAjVsvQY7AGc7xZVFLjPc1Z";
	private String password = "";
	
	private static final String gasFee = "100";
	private static final String COIN_XRP = "XRP";
 
    private Logger log = Log.get();
 
    private final static String RESULT = "result";
    private final static String SUCCESS = "success";
    private final static String TES_SUCCESS = "tesSUCCESS";
    
    private final static String METHOD_GET_TRANSACTION = "/v2/accounts/{0}/transactions";
    private final static String METHOD_GET_BALANCE = "/v2/accounts/{0}/balances";
    
    private final static String METHOD_POST_INDEX = "ledger_current";
    private final static String METHOD_POST_ACCOUNT_INFO = "account_info";
    private final static String METHOD_POST_SUBMIT = "submit";
    
    public static void main(String args[]) throws Exception{
    	CopyOfCoinXrpService XRPUtils = new CopyOfCoinXrpService();
//    	System.out.println(XRPUtils.send("rG2pkYL2q9jDnEA7xraKH7gXLXo17Tj4tX", 0.1));
//    	System.out.println(XRPUtils.parseTransaction("0"));
    	System.out.println(XRPUtils.sign("rG2pkYL2q9jDnEA7xraKH7gXLXo17Tj4tX",1000.0));
    }
    
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
    	value = BigDecimalUtil.mul(value, 1000000);
    	Integer vInteger = BigDecimal.valueOf(value).intValue();
    	Map<String, String> map = getAccountSequenceAndLedgerCurrentIndex();
    	Payment payment = new Payment();
        payment.as(AccountID.Account,     address);
        payment.as(AccountID.Destination, toAddress);
        payment.as(UInt32.DestinationTag, "1");
        payment.as(Amount.Amount,         vInteger.toString());
        payment.as(UInt32.Sequence,       map.get("accountSequence"));
        payment.as(UInt32.LastLedgerSequence, map.get("ledgerCurrentIndex")+4);
        payment.as(Amount.Fee,            gasFee);
        SignedTransaction signed = payment.sign(password);
        if (signed != null) {
        	return signed.tx_blob;
		}
        return null;
    }
 
    public Map<String, String> getAccountSequenceAndLedgerCurrentIndex(){
    	HashMap<String, String> params = new HashMap<String, String>();
    	params.put("account", address);
    	params.put("strict", "true");
    	params.put("ledger_index", "current");
    	params.put("queue", "true");
    	JSONObject re = doRequest(METHOD_POST_ACCOUNT_INFO,params);
        if (re != null) {
        	JSONObject  result = re.getJSONObject("result");
			if (SUCCESS.equals(result.getString("status"))) {
				Map<String, String> map = new HashMap<String, String>();
				map.put("accountSequence", result.getJSONObject("account_data").getString("Sequence"));
				map.put("ledgerCurrentIndex", result.getString("ledger_current_index"));
				return map;
			}
		}
        return null;
    }
    
    /**
     * 获取用户交易序列号
     * @return
     */
    public long getAccountSequence(){
    	HashMap<String, String> params = new HashMap<String, String>();
    	params.put("account", address);
    	params.put("strict", "true");
    	params.put("ledger_index", "current");
    	params.put("queue", "true");
    	JSONObject re = doRequest(METHOD_POST_ACCOUNT_INFO,params);
        if (re != null) {
        	JSONObject  result = re.getJSONObject("result");
			if (SUCCESS.equals(result.getString("status"))) {
				return result.getJSONObject("account_data").getLongValue("Sequence");
			}
		}
        return 0L;
    }
    
    /**
     * 获取最新序列
     * @return
     */
    public long getLedgerIndex(){
    	JSONObject re = doRequest(METHOD_POST_INDEX);
        if (re != null) {
        	JSONObject  result = re.getJSONObject("result");
			if (SUCCESS.equals(result.getString("status"))) {
				return result.getLongValue("ledger_current_index");
			}
		}
        return 0L;
    }
    
    /**
     * XRP查询余额
     * @return
     */
    public double getBalance(){
    	HashMap<String, String> params = new HashMap<String, String>();
    	params.put("currency", COIN_XRP);
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
    		Date d = new Date(BigDecimalUtil.longAdd(Long.parseLong(startTm), 1000L));
        	params.put("start", DateUtil.dateToString(d, "yyyy-MM-dd'T'HH:mm:ss"));
		}
    	params.put("result", "tesSUCCESS");
    	params.put("type", "Payment");
        //如若考虑1秒并发100条以上记录.请不要用时间分页.用marker读取账户交易记录,或者ledger_index进度来处理.
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
						//校验用户是否存在
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
						//添加充值记录
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
--------------------- 
作者：其修远兮 
来源：CSDN 
原文：https://blog.csdn.net/liu1765686161/article/details/83347534 
版权声明：本文为博主原创文章，转载请附上博文链接！