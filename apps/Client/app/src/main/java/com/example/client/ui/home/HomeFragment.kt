package com.example.client.ui.home

import android.app.Activity
import android.content.ContentResolver
import android.content.ContentUris
import android.content.Context
import android.content.Intent
import android.database.Cursor
import android.media.MediaMetadataRetriever
import android.net.Uri
import android.os.Bundle
import android.os.Environment
import android.provider.DocumentsContract
import android.provider.MediaStore
import android.provider.OpenableColumns
import android.util.Log
import android.view.View
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.fragment.app.Fragment
import by.kirich1409.viewbindingdelegate.viewBinding
import com.example.client.Constants.CHUNK_DURATION
import com.example.client.Constants.LOG_TAG_MEDIA
import com.example.client.Constants.LOG_TAG_UI
import com.example.client.databinding.FragmentHomeBinding
import com.example.client.R
import com.example.client.media.splitVideo
import java.io.File

class HomeFragment : Fragment(R.layout.fragment_home) {

    private val b by viewBinding(FragmentHomeBinding::bind)

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        setClickListeners()
    }

    fun setClickListeners() {
        b.clickMe.setOnClickListener {
            openFileSelector()
        }
    }

    private val pickVideoLauncher =
        registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
            if (result.resultCode == Activity.RESULT_OK) {
                val data: Intent? = result.data
                val videoUri = data?.data ?: return@registerForActivityResult

                val fileName = getFileName(requireContext().contentResolver, videoUri)
//                val video = Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_MOVIES + File.separator + fileName)
                val video = File("/Internal storage/Android/media/com.whatsapp/WhatsApp/Media/WhatsApp Video/$fileName")
                if (video.exists()) {
                    splitVideo(video, CHUNK_DURATION)
                }
            }
        }

    fun getFileName(contentResolver: ContentResolver, uri: Uri): String {
        var fileName = ""
        val cursor = contentResolver.query(uri, null, null, null, null)
        cursor?.use {
            if (it.moveToFirst()) {
                val columnIndex = it.getColumnIndex(OpenableColumns.DISPLAY_NAME)
                if (columnIndex > -1) {
                    fileName = it.getString(columnIndex)
                    Log.d(LOG_TAG_MEDIA, "File Name: $fileName")
                }
            }
        }
        return fileName
    }


    private fun getFilePathFromURI(context: Context, contentUri: Uri): String? {
        val resolver = context.contentResolver
        val fileName = getFileName(contentUri, resolver) ?: return null
        val outputFile = File(context.cacheDir, fileName)
        resolver.openInputStream(contentUri)?.use { input ->
            outputFile.outputStream().use { output ->
                input.copyTo(output)
            }
        }
        return outputFile.absolutePath
    }

    private fun getFileName(uri: Uri, resolver: ContentResolver): String? {
        val cursor = resolver.query(uri, null, null, null, null)
        val name = if (cursor != null && cursor.moveToFirst()) {
            val index = cursor.getColumnIndex(OpenableColumns.DISPLAY_NAME)
            if (index != -1) {
                cursor.getString(index)
            } else null
        } else null
        cursor?.close()
        return name
    }

    private fun openFileSelector() {
        val intent = Intent(Intent.ACTION_OPEN_DOCUMENT)
        intent.addCategory(Intent.CATEGORY_OPENABLE)
        intent.type = "video/*" // Allow only video files
        pickVideoLauncher.launch(intent)
    }
}