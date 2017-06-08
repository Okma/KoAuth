package com.carlamko.apps.koauth.ui

import android.app.AlertDialog
import android.app.Dialog
import android.app.DialogFragment
import android.os.Bundle
import android.view.View
import com.carlamko.apps.koauth.R
import android.content.Intent
import android.app.Activity
import android.content.SharedPreferences
import android.preference.PreferenceManager
import android.text.Editable
import android.text.TextWatcher
import android.widget.Button
import kotlinx.android.synthetic.main.add_serial_dialog.view.*

/**
 * Created by Carl on 5/31/2017.
 */
class AddSerialDialog : DialogFragment() {

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val dialogView: View = activity.layoutInflater.inflate(R.layout.add_serial_dialog, null)
        val builder: AlertDialog.Builder = AlertDialog.Builder(activity)

        builder.setView(dialogView)
                .setPositiveButton(R.string.dialog_positive_text, {
                    dialog, _ ->

                    val editor: SharedPreferences.Editor = PreferenceManager.getDefaultSharedPreferences(activity).edit()
                    editor.putString(getString(R.string.pref_serial_key), dialogView.et_serial.text.toString())
                    editor.apply()

                    targetFragment.onActivityResult(targetRequestCode, Activity.RESULT_OK, Intent())
                })
                .setNegativeButton(R.string.dialog_negative_text, {
                    dialog, _ ->
                    dialog.cancel()
                })

        val alertDialog: AlertDialog = builder.create()

        alertDialog.setOnShowListener {
            // Turn off positive button.
            val positiveButton: Button = (dialog as AlertDialog).getButton(AlertDialog.BUTTON_POSITIVE)
            positiveButton.isEnabled = positiveButton.text.isNotEmpty()

            dialogView.et_serial.addTextChangedListener(object : TextWatcher {
                override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}
                override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
                override fun afterTextChanged(s: Editable?) {
                    if (s != null) {
                        positiveButton.isEnabled = s.isNotEmpty()
                    }
                }
            })
        }

        return alertDialog
    }
}