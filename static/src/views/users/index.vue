<template>
  <el-table
    :data="tableData"
    stripe
    style="width: 100%">
    <el-table-column
      prop="id"
      label="用户编号"
      width="80">
    </el-table-column>
    <el-table-column
      prop="reg_date"
      :formatter="formatDateTime"
      label="注册日期"
      width="100">
    </el-table-column>
 <el-table-column
      prop="pay_date"
      label="服务开始启用日期"
      width="200">
    </el-table-column>
    

    <el-table-column
      prop="username"
      label="姓名"
      width="100">
    </el-table-column>
     <el-table-column
      prop="password"
      label="密码"
      width="100">
    </el-table-column>
    <el-table-column
      prop="phone"
       width="120"
      label="电话号码">
    </el-table-column>
     <el-table-column
      prop="pay_type"
       width="120"
      :formatter="formatPayType"
      label="付费服务类型">
    </el-table-column>

     <el-table-column
      
      label="操作">
      <template slot-scope="scope">
        <el-button
          size="mini"
          @click="handleUpdateNormal(scope.$index, scope.row)">升级普通会员</el-button>
        <el-button
          size="mini"
          @click="handleUpdateVip(scope.$index, scope.row)">升级高级会员</el-button>
        <el-button
          size="mini"
          type="danger"
          @click="handleDelete(scope.$index, scope.row)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
import { getUsers } from '@/api/login'
import { updateUser } from '@/api/login'
import { deleteUser } from '@/api/login'

export default {
  data() {
    return {
      tableData: []
    }
  },
  methods: {
    handleUpdateNormal(index, row) {
      console.log(row)
      var user = {
        pay_type: 1
      }
      updateUser(row.id, user).then(response => {
        this.fetchData()
      })
    },
    handleUpdateVip(index, row) {
      var user = {
        pay_type: 2
      }
      updateUser(row.id, user).then(response => {
        this.fetchData()
      })
    },
    handleDelete(index, row) {
      deleteUser(row.id).then(response => {
        this.fetchData()
      })
    },
    fetchData() {
      this.listLoading = true
      getUsers().then(response => {
        this.tableData = response.data
        this.listLoading = false
      })
    },
    formatDatetwo: function(time) {
      var re = /-?\d+/
      var m = re.exec(time)
      var d = new Date(parseInt(m[0]))
      var o = {
        'M+': d.getMonth() + 1, // month
        'd+': d.getDate(), // day
        'h+': d.getHours(), // hour
        'm+': d.getMinutes(), // minute
        's+': d.getSeconds(), // second
        'q+': Math.floor((d.getMonth() + 3) / 3), // quarter
        'S': d.getMilliseconds() // millisecond
      }
      var format = 'yyyy-MM-dd'
      if (/(y+)/.test(format)) {
        format = format.replace(RegExp.$1, (d.getFullYear() + '').substr(4 - RegExp.$1.length))
      }
      for (var k in o) {
        if (new RegExp('(' + k + ')').test(format)) {
          format = format.replace(RegExp.$1, RegExp.$1.length === 1 ? o[k] : ('00' + o[k]).substr(('' + o[k]).length))
        }
      }
      return format
    },
    formatDateTime(row, column, cellValue, index) {
      return this.formatDatetwo(row.reg_date * 1000)
    },
    formatPayType(row, column, cellValue, index) {
      console.log(row, column, cellValue, index)

      if (row.pay_type === 0) {
        return '普通用户'
      } else if (row.pay_type === 1) {
        return '普通付费用户'
      } else if (row.pay_type === 2) {
        return '高级付费用户'
      }
    }
  },
  mounted() {
    console.log('mouted')
    this.fetchData()
  },
  filters: {

  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

