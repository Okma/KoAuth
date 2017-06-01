package com.carlamko.apps.koauth

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.v7.widget.Toolbar
import android.view.Menu
import android.view.MenuItem

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        setSupportActionBar(findViewById(R.id.toolbar_main) as Toolbar)
        supportActionBar!!.setDisplayShowTitleEnabled(false)
    }

    override fun onCreateOptionsMenu(menu: Menu): Boolean {
        menuInflater.inflate(R.menu.settings_menu, menu)
        return true
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        // Settings option selected.
        if (item.itemId == R.id.settings_icon) {
            val openSettings = Intent(this, SettingsActivity::class.java)
            startActivity(openSettings)
        }
        return false
    }
}