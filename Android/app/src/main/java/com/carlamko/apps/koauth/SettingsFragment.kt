package com.carlamko.apps.koauth

import android.app.Fragment
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import com.carlamko.apps.koauth.ui.AddSerialDialog

class SettingsFragment : Fragment(), View.OnClickListener {

    private val FRAGMENT_ID :Int = 1

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
    }

    override fun onCreateView(inflater: LayoutInflater?, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        val view: View = inflater!!.inflate(R.layout.fragment_settings, container, false)

        //view.findViewById(R.id.tv_serial)

        return view
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
