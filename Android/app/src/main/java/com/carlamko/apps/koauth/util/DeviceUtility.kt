package com.carlamko.apps.koauth.util

import android.os.Build

/**
 * Created by Carl on 5/30/2017.
 */

object DeviceUtility {

    private val TEST_SERIAL: String = "testserial123"

    fun getDeviceSerial(): String {
        if(!Build.SERIAL.isEmpty() && Build.SERIAL != "unknown") {
            return Build.SERIAL
        } else {
            return TEST_SERIAL
        }
    }
}
