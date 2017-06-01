package com.carlamko.apps.koauth.ui

import android.app.AlertDialog
import android.app.Dialog
import android.app.DialogFragment
import android.content.DialogInterface
import android.os.Bundle
import android.view.View
import com.carlamko.apps.koauth.R

/**
 * Created by Carl on 5/31/2017.
 */
class AddSerialDialog : DialogFragment() {

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val dialogView: View = activity.layoutInflater.inflate(R.layout.add_serial_dialog, null)
        val builder: AlertDialog.Builder = AlertDialog.Builder(activity);

        builder.setView(dialogView)
                .setPositiveButton(R.string.dialog_positive_text, DialogInterface.OnClickListener {
                    dialog, which ->

                    activity
                })
                .setNegativeButton(R.string.dialog_negative_text, DialogInterface.OnClickListener {
                    dialog, which ->
                    dialog.cancel()
                })

        return super.onCreateDialog(savedInstanceState)
    }
}