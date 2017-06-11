package com.carlamko.apps.koauth

import android.app.Fragment
import android.content.Intent
import android.content.SharedPreferences
import android.os.Bundle
import android.preference.PreferenceManager
import android.support.design.widget.Snackbar
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_settings.view.*
import com.carlamko.apps.koauth.ui.AddSerialDialog
import com.carlamko.apps.koauth.util.NetworkUtility

class SettingsFragment : Fragment(), View.OnClickListener {
    private val FRAGMENT_ID :Int = 1

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
    }

    override fun onCreateView(inflater: LayoutInflater?, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        val view: View = inflater!!.inflate(R.layout.fragment_settings, container, false)

        val sharedPreferences: SharedPreferences = PreferenceManager.getDefaultSharedPreferences(activity)
        val currentSerial = sharedPreferences.getString(getString(R.string.pref_serial_key), null)
        if(currentSerial != null) {
            view.tv_serial.text = currentSerial
            view.fab_add_serial.visibility = View.GONE
        } else {
            view.fab_add_serial.setOnClickListener(this)
        }
        return view
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        when(requestCode) {
            FRAGMENT_ID -> {
                val sharedPreferences: SharedPreferences = PreferenceManager.getDefaultSharedPreferences(activity)

                // Fetch the new serial value.
                val serialInput: String? = sharedPreferences.getString(getString(R.string.pref_serial_key), null)

                if(serialInput != null) {
                    // Check that new serial value is valid.
                    NetworkUtility.checkSerial(activity, serialInput,
                            {
                                _ ->
                                view.tv_serial.text = String.format(getString(R.string.serial_format), serialInput)
                                view.fab_add_serial.visibility = View.GONE

                                Snackbar.make(view, "Serial verified!", Snackbar.LENGTH_LONG).show()
                            },
                            {
                                error ->
                                view.tv_serial.text = "Error: ${error.message}"

                                Snackbar.make(view, "Invalid serial!", Snackbar.LENGTH_LONG).show()
                            })
                }
            }
        }
    }

    override fun onClick(v: View?) {
        when(v?.id) {
            R.id.fab_add_serial -> {
                val newDialog: AddSerialDialog = AddSerialDialog()
                newDialog.setTargetFragment(this, FRAGMENT_ID)
                newDialog.show(fragmentManager, "add-serial")
            }
        }
    }

}
