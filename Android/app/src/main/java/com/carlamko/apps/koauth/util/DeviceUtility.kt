package com.carlamko.apps.koauth.util

import android.os.Build

/**
 * Created by Carl on 5/30/2017.
 */

object DeviceUtility {

    val deviceSerial: String
        get() = Build.SERIAL
}
