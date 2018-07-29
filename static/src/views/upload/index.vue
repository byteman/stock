<template>
  <div class="app-container">
  <el-upload
  class="upload-demo"
  :action="postfile"
  :on-preview="handlePreview"
  :on-success="handleOnSuccess"
  :file-list="fileList">
  <el-button size="small" type="primary">点击上传</el-button>
  <div slot="tip" class="el-upload__tip">只能上传txt文件</div>
</el-upload>


  </div>
</template>

<script>
  export default {
    data() {
      return {
        fileList: [],
        postfile: process.env.BASE_API + '/upload'
      }
    },
    method: {
      handleRemove(file, fileList) {
        console.log(file, fileList)
      },
      handleOnSuccess(response, file, fileList) {
        this.$message({
          message: '文件上传成功',
          center: true
        })
      }
  
    }
  }
  
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

