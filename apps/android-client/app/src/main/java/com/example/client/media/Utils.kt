package com.example.client.media

import android.content.Context
import android.media.MediaMetadataRetriever
import android.net.Uri
import android.os.Environment
import android.util.Log
import androidx.annotation.Nullable
import java.io.File
import com.arthenica.mobileffmpeg.FFmpeg
import com.arthenica.mobileffmpeg.Config
import com.arthenica.mobileffmpeg.FFprobe
import com.example.client.Constants.LOG_TAG_MEDIA
import java.io.FileDescriptor
import java.io.FileOutputStream
import java.util.concurrent.Executors

class MediaError(code: String, message: String) : Throwable("${code}_error: $message")

fun splitVideo(videoFile: File, chunkDuration: Int = 10): Throwable? {

    val outputDir = File(Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_MOVIES), "VideoChunks")
    outputDir.mkdirs()
    val outputFile = File(outputDir, "chunk_0.mp4")

    val command = "-ss 0 -t $chunkDuration -i ${videoFile.absolutePath} -c copy ${outputFile.absolutePath}"
    val returnCode = FFmpeg.execute(command)

//    val duration = FFprobe.getMediaInformation(inputFilePath).duration.toLong() * 1000
//
//    var startTime: Long = 0
//    var chunkIndex = 0
//    while (startTime < duration) {
//        val outputFile = File(outputDir, "chunk_$chunkIndex.mp4")
//        val endTime = startTime + (chunkDuration * 1000)
//        if (endTime > duration) {
//            break
//        }
//
//        val command = "-ss $startTime -t $chunkDuration -i $inputFilePath -c copy ${outputFile.absolutePath}"
//        val returnCode = FFmpeg.execute(command)
//
//        if (returnCode == Config.RETURN_CODE_SUCCESS) {
//            startTime = endTime
//            chunkIndex++
//            Log.d(LOG_TAG_MEDIA, "Chunk $chunkIndex created.\nOutput File: ${outputFile.absolutePath}")
//        } else {
//            // TODO: handle error
//            Log.e(LOG_TAG_MEDIA, "Return code: $returnCode.\nFailed to execute command: $command.")
//            return MediaError("ffmpeg", "Failed to split video.")
//        }
//    }

    return null
}