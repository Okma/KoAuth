package com.carlamko.apps.koauth.util

import android.content.Context
import android.net.Uri
import com.android.volley.Request
import com.android.volley.RequestQueue
import com.android.volley.Response
import com.android.volley.VolleyError
import com.android.volley.toolbox.JsonObjectRequest
import com.android.volley.toolbox.Volley
import org.json.JSONObject

/**
 * Created by Carl on 6/7/2017.
 */

object NetworkUtility {

    private var requestQueue: RequestQueue? = null
    private val uriBuilder: Uri.Builder = Uri.Builder()
    private val CHECK_SERIAL_URL: String by lazy {
        uriBuilder.scheme("http").encodedAuthority("10.0.2.2:8080").appendPath("validate").appendPath("serial").build().toString()
    }

    private fun getRequestQueue(context: Context): RequestQueue {
        if(requestQueue == null) {
            requestQueue = Volley.newRequestQueue(context)
        }
        return requestQueue as RequestQueue
    }


    fun checkSerial(context: Context, serial: String, successCallback: (JSONObject) -> Unit, failureCallback: (VolleyError) -> Unit) {
        val json: JSONObject = JSONObject()
        json.put("serial", serial)
        json.put("device_serial", DeviceUtility.deviceSerial)

        getRequestQueue(context).add(object: JsonObjectRequest(Request.Method.POST, CHECK_SERIAL_URL,
                json,
        Response.Listener<JSONObject> {
            response ->  successCallback(response) },
        Response.ErrorListener {
            error ->  failureCallback(error) }) {})
    }

}