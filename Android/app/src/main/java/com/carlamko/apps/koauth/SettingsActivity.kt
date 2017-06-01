package com.carlamko.apps.koauth

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.v7.widget.Toolbar
import android.view.MenuItem

class SettingsActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_settings)

        setSupportActionBar(findViewById(R.id.toolbar_settings) as Toolbar)
        supportActionBar!!.setDisplayHomeAsUpEnabled(true)
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        if (item.itemId == android.R.id.home) {
            onNavigateUp()
        }

        return super.onOptionsItemSelected(item)
    }
}
