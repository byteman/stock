/**
 * Created by lz on 2017/6/28.
 */
var obj = new Api();
var myflage = 0;
var cluflag = 1;
var txtflag = 1;
var levelflag;
var pageSize = 20;//每页显示数据
var pageNo = '1';//当前页数
var pageNo2 = '1';
var la;
var lg;
var jgid;
//org_id组织编号
var  org_id;
var boundary;
var nowCenter, nowCenterLng, nowCenterLat;
var markers = [];
var labels = [];
//winHeight:屏幕窗口高度，winWidth：屏幕窗口宽度
var winHeight = $(window).height();
// var winWidth = $(window).width();
//topHeight顶部高度
var topHeight = 130;
//panelHeight 面板总高度,绘画区域高度
var panelHeight = winHeight - topHeight;
var marker2=[];
// var map = new BMap.Map("allmap",{
//     minZoom : 4,
//     maxZoom : 16,
//     mapType:BMAP_SATELLITE_MAP
// });
//var jgid=getUrlParam("jgid");   //组织编号
var mapSource = "BD"; 			  //地图源 add by lwh at 2017-12-13
var iconPng = ["icon-green.png","icon-orang.png","icon-hs.png"]; //地图标注使用图片 add by lwh at 2017-12-13
map = new BMap.Map("allmap",{mapType:BMAP_NORMAL_MAP});
var markerClusterer = new BMapLib.MarkerClusterer();
var refreshQmapFlag = false;

/* 腾讯地图简明函数申明开始 add by lwh at 2017-12-15 */
var Map = qq.maps.Map;
var Marker = qq.maps.Marker;
var LatLng = qq.maps.LatLng;
var Event = qq.maps.event;
var MarkerImage = qq.maps.MarkerImage;
var MarkerShape = qq.maps.MarkerShape;
var MarkerAnimation = qq.maps.MarkerAnimation;
var Point = qq.maps.Point;
var Size = qq.maps.Size;
var ALIGN = qq.maps.ALIGN;
var MVCArray = qq.maps.MVCArray;
var MarkerCluster = qq.maps.MarkerCluster;
var Cluster = qq.maps.Cluster;
var MarkerDecoration = qq.maps.MarkerDecoration;
//清除覆盖物的函数  eg:clearQmapOverlays(markers);
function clearQmapOverlays(overlays){
	var overlay;
	while(overlay = overlays.pop())
	{
		overlay.setMap(null);
	}
}
//关闭聚合
function closeQmapClusterMarkers() {
	//markerClusterer.clearMarkers(markers);
	return;
}
//创建聚合
function clusterQmapMarkers() {
	markerClusterer = new MarkerCluster({
		map:map,
		minimumClusterSize:2, //默认2
		markers:markers,
		zoomOnClick:true, //默认为true
		gridSize:60, //默认60
		averageCenter:false, //默认false
		maxZoom:18, //默认18
	});
	Event.addListener(markerClusterer, 'clusterclick', function (evt) {
		var ss =  evt;
	});
}
//添加标注
function addQmapMarker(markers, labels) {
    var i = 0;
	for (; i < markers.length; i++) {
		var mymarkadd = markers[i];		
		mymarkadd.setMap(map); // 添加单个覆盖物
		Event.addListener(mymarkadd, 'click', function() {
			location.href = "cropCon.html?id="+this.id + "&orgId=" + this.orgId;
		});
		//必要时可给Label添加click事件
	}	
    if(txtflag){
        showLabel(labels);
        myflage = 1;
    }else if(!txtflag){
        removeLabel(labels);
        myflage = 1;
    }
}

//转换GPS坐标为腾讯坐标
var tmpGlbQQLatLng = new qq.maps.LatLng(0, 0);
function GPS2QQ(lat,lng){
	qq.maps.convertor.translate(new qq.maps.LatLng(lat, lng), 1, function(res){tmpGlbQQLatLng = res[0];map.panTo(tmpGlbQQLatLng)});
}
/* 腾讯地图简化函数申明结束 */

//
function findOldMap()
{
	var curSel = $("#maptype").val();
	var oldSelMap = getMainSelectMap();
	if((oldSelMap == null) || (curSel == oldSelMap))
		return 0; // 未发生变化
	if(refreshQmapFlag == true)
		return 0; // QQ地图需要刷新
	$("#maptype").val(oldSelMap);
	$('#maptype').trigger("change");
	return 1; // 发生了变化
}

// 获取登录信息
obj.get('/api/v1/users', '', function(data) {
    for (var i = 0; i < data.length; i++) {
        var data = data[i];
        $("#name").html(data.username);
        org_id = data.org_id;
        $("#href_stationStatus").attr("href","stationStatus.html?orgId="+org_id);
	    $("#href_dataSearch").attr("href","dataSearch.html?orgId="+org_id);
        $("#href_analyse").attr("href","analyse.html?orgId="+org_id);
        $("#href_corpBasicInfo").attr("href","corps.html?orgId="+org_id);
    }
}, function(err) {
    $("#name").html("");
});
//获取登录用户行政区域
obj.get('/api/v1/gis/boundary', '', function(data) {
    boundary = data.boundary;
    la = data.latitude;
    lg = data.longitude;
    var account_type = data.account_type;
    levelflag = account_type;
    if(levelflag==2){
        $("#myfocusLi").hide();
        $("#focus-list").hide();
        $(".messagelist_title li").css('width','50%');
        // map.centerAndZoom(new BMap.Point(lg,la), 11);
        // map.enableScrollWheelZoom();
        // map.panTo(new BMap.Point(lg+0.001,la));
    }else if(levelflag==0 || levelflag == 1){
        // map.centerAndZoom(new BMap.Point(lg,la), 10);
        // map.enableScrollWheelZoom();
        // map.panTo(new BMap.Point(lg+0.001,la));
        $("#myfocusLi").show();
    }
    $('#levelid').val('');
    getAllMarkers();
    getTree();
    getMessageList(jgid,pageSize,pageNo);
    getWarnList(bef_data,now_data,pageSize,pageNo2);
    getFocus();
    getWeather(boundary);
}, function(err) {
    return;
});
$("#allmap").height(panelHeight);
$("#messageBox").height(panelHeight-190);
$(".tableheight").height($('#messageBox').height()-50);
$(".tableheight2").height($('#messageBox').height()-83);
$("#treeDemo").css("max-height",panelHeight-362+"px");
$("#href_mainPage").addClass('curOptMenu');
//控制左侧面板是否显示
$("#par_control").click(function() {
    if ($("#par_showIf").css("left") == "-220px") {
        $("#par_showIf").animate({ left: 0, opacity: 1 }, 500);
        $("#par_control").removeClass("par_control_left").attr("title", "隐藏面板");
        $("#dragZoomControl").animate({ left: "220px" }, 500);
    } else {
        $("#par_showIf").animate({ left: "-220px", opacity: 0 }, 500);
        $("#par_control").addClass("par_control_left").attr("title", "显示面板");
        $("#dragZoomControl").animate({ left: "42px" }, 500);
    }
});
function reloadMap()
{
    var lg= map.getCenter().lng;
    var la =map.getCenter().lat;
    var zoom = map.getZoom();
	
    getAllMarkers();	
	
}
$("#maptype").change(function () {
    var type = $("#maptype").val();
    if(type == "百度-卫星地图")
    {
        mapSource = "BD";
		setMainSelectMap("百度-卫星地图");
		map = new BMap.Map("allmap",{mapType:BMAP_HYBRID_MAP});
		markerClusterer = new BMapLib.MarkerClusterer();
		reloadMap();
    }
    else if(type == "百度-街道地图")
    {
        mapSource = "BD";
		setMainSelectMap("百度-街道地图");
		map = new BMap.Map("allmap",{mapType:BMAP_NORMAL_MAP});
		markerClusterer = new BMapLib.MarkerClusterer();
		reloadMap();
    }
    else if(type == "腾讯-卫星地图") // QQ_HYBRID
	{
		mapSource = "QQ";
		setMainSelectMap("腾讯-卫星地图");
		map = new qq.maps.Map(document.getElementById("allmap"), {			
			zoom: 10,							// 地图默认缩放等级
			center: new qq.maps.LatLng(la, lg), // 地图的中心地理坐标			
			 mapTypeId: qq.maps.MapTypeId.HYBRID //该地图类型显示卫星图像上的主要街道透明层，也可使用SATELLITE但无街道名称
		});
		if(cluflag)
			clusterQmapMarkers();		
		reloadMap();
    }
	else if(type == "腾讯-街道地图") // QQ_NORMAL
	{
		mapSource = "QQ";
		setMainSelectMap("腾讯-街道地图");
		map = new qq.maps.Map(document.getElementById("allmap"), {			
			zoom: 10,							// 地图默认缩放等级
			center: new qq.maps.LatLng(la, lg), // 地图的中心地理坐标			
			 mapTypeId: qq.maps.MapTypeId.ROADMAP //该地图类型显示普通的街道地图
		});
		if(cluflag)
			clusterQmapMarkers();
		reloadMap();
    }
	else
	{
		;//do nothing
	}
})
//控制右侧面板是否显示
$("#par_control_ri").click(function() {
    if ($("#par_showRt").css("right") == "-458px") {
        $("#par_showRt").animate({ right: 0, opacity: 1 }, 500);
        $("#ppar_control_ri").addClass("par_control-ri-left").attr("title", "隐藏面板");
        $("#dragZoomControl-ri").animate({ right: "458px" }, 500);

    } else {
        $("#par_showRt").animate({ right: "-458px", opacity: 0 }, 500);
        $("#par_control_ri").addClass("par_control-ri-left").attr("title", "显示面板");
        $("#dragZoomControl-ri").animate({ right: "-27px" }, 500);
    }
});

//当前时间
var now_day=new Date().Format("yyyy-MM-dd hh:mm");
//一个月前
var bef_month=getDateForMonthAlong(now_day,1);
var now_data=parseInt(new Date(now_day).getTime()/1000);
var bef_data=parseInt(new Date(bef_month).getTime()/1000);
$("#starttime").val(bef_month);
$("#endtime").val(now_day);

//设置预警列表时间段
$('#starttime').datetimepicker({
    language: "zh-CN",
    format: 'yyyy-mm-dd hh:ii',
    autoclose: true,
    todayBtn: true,
})
$('#endtime').datetimepicker({
    language: "zh-CN",
    format: 'yyyy-mm-dd hh:ii',
    autoclose: true,
    todayBtn: true,
}) ;


//文件树操作初始设置
var setting = {
    check: {
        enable: true,
        chkStyle: "radio",
        radioType: "all"
    },
    data: {
        simpleData: {
            enable: true,
            pIdKey: "pid"
        }
    },
    view: {
        txtSelectedEnable: false,
        addDiyDom: addDiyDom
    },
    callback: {
        onClick: zTreeOnClick,
        onCheck: zTreeOnClick,
        onAsyncSuccess: zTreeOnAsyncSuccess
    }
};
// function getTree(){
//     //获取文件树结构
//     obj.get('/api/v1/orgs', '', function(data) {
//         $.fn.zTree.init($("#treeDemo"), setting, data.orgs);
//         var treeObj = $.fn.zTree.getZTreeObj("treeDemo");
//         $("#treeDemo_1_check").addClass('radio_true_full');
//         treeObj.expandAll(true);
//     });
// }
function logout()
{
    window.localStorage.removeItem('username');
    window.localStorage.removeItem('password');
    window.location.href="./index.html";

}
function SelectNode() {
    jgid = getMainCorpSelectId();//getUrlParam("jgid")
    //console.log("jgid=====",jgid)
    // if(jgid===undefined)
    // {return;}
    // jgid=undefined
    var treeObj = $.fn.zTree.getZTreeObj("treeDemo");
    if(treeObj===null)
    {
        return;
    }
    var treeNode = treeObj.getNodeByParam("id",jgid,null)
    if(treeNode===null)
    {
        ;//console.log("can not find jgid" ,jgid);
    }else
    {
        ;//console.log("find jgid="+jgid+"name=",treeNode.name);
        zTreeOnClick(null,jgid,treeNode);
    }
}
//获取尾矿库信息列表
function getMessageList(jgid,pageSize,pageNo) {
    var data_m = {
        org_id: jgid,
        page_size:pageSize,
        page: pageNo
    };
    //console.log("getMessageList jgid=",jgid)
    obj.get('/api/v1/corps',data_m,function(data){
        var mypageNo = data.current_page==null?0:data.current_page;
        var mypages = data.total_pages==null?0:data.total_pages;
        $("#pageMsg").html(mypageNo+"/"+mypages);
        $("#curpage").val(mypageNo);
        $("#tolpage").val(mypages);
        $('#toPage').val(mypageNo);
        $("#talOne").html(data.total_records);
        var mydata = data.data;
        $("#tableBody").html("");
        for(var i= 0;i<mydata.length;i++){
            var object = mydata[i];
            showMessage(object);
        }

    },function(err){
        return;
    })
}

//尾矿库列表信息展示
function showMessage(object){
    var html = $("#tableBody").html();
    var trHTML = "<tr class=\"tbody-tr\" onclick=\"gotoMap('"+object.latitude+"','"+object.longitude+"')\">";
    trHTML += '<td class="tbody-td">'+object.id+'</td>';
    trHTML += '<td class="tbody-td"><a class="tbody-td-a" href="cropCon.html?id='+object.id+'&orgId='+object.my_org_id+'">'+object.corp_name+'</a></td>';
    
	trHTML += "<td class=\"tbody-td\">"+object.region+"</td>";
    $("#tableBody").html(html+trHTML);
}
function gotoMap(la,lg){
	if(mapSource == "BD")
	{
		map.centerAndZoom(new BMap.Point(lg,la),16);  
		marker2 = new BMap.Marker(new BMap.Point(lg,la));
	}
	else if(mapSource == "QQ") //add by lwh at 2017-12-21 
	{
		//GPS2QQ(la,lg);//
		map.panTo(new qq.maps.LatLng(la, lg));
        map.setZoom(16);
	}
}
function gotoMap2(la,lg,jgid,treeflag,name){
    if(mapSource == "BD")
		nowCenter = new BMap.Point(lg, la);
	else if(mapSource == "QQ")
		nowCenter = new qq.maps.LatLng(la, lg);
    if(treeflag && name!=boundary){//父级
        map.panTo(nowCenter);
        map.setZoom(14);
    }else if(!treeflag){//子级
        map.panTo(nowCenter);
        map.setZoom(16);
    }
	if(mapSource == "QQ") // add by lwh at 2017-12-19
	{
		//console.log("sssss lg"+lg+"lat:"+la);// ???????????????????
		return;
	}		
    if(name == boundary){
        getAllMarkers();
        if(levelflag==2){
            map.setZoom(11);
        }else if(levelflag==0 || levelflag == 1){
            map.setZoom(10);
        }
    }
}
//获取预警列表
function getWarnList(starttime,endtime,pageSize,pageNo2,level){
    var data_m = {
        start:starttime,
        end:endtime,
        page_size:pageSize,
        page: pageNo2,
        level : level
    };
    obj.get('/api/v1/alarm',data_m,function(data){
        var mypageNo2 = data.current_page==null?0:data.current_page;
        var mypages2 = data.total_pages==null?0:data.total_pages;
        $("#pageMsg2").html(mypageNo2+"/"+mypages2);
        $("#curpage2").val(mypageNo2);
        $("#tolpage2").val(mypages2);
        $('#toPage2').val(mypageNo2);
        $("#talTwo").html(data.total_records);
        var _data=data.data;
        $("#warningBody").html("");
        $("#bjNum").html('');
        $("#bjTime").html('')
        for(var i=0;i<_data.length;i++){
            var object=_data[i];
            showWarnMessage(object);
        }
    },function(err){
        return;
    })
}
//预警列表展示
function showWarnMessage(object){
    var html = $("#warningBody").html();
    var trHTML = "<tr class=\"tbody-tr\">";
    trHTML += "<td class=\"tbody-td\"><label class=\"mydelbtn\" title=\"处理预警\" onclick=\"takeWarn('"+object.id+"','"+getMyDate(object.time)+"','"+object.value+"','"+object.status+"','"+object.user+"','"+object.alarm_type+"','"+object.detail+"')\"></label>"+object.alarm_type+"</td>";
    trHTML += "<td class=\"tbody-td\">"+object.station.name+"</td>";
    trHTML += "<td class=\"tbody-td\">"+object.alarm.type+"</td>";
    trHTML += "<td class=\"tbody-td\">"+object.station.name+"</td>";
    trHTML += "<td class=\"tbody-td\">"+getMyDate(object.time)+"</td>";
    trHTML += "</tr>";
    $("#warningBody").html(html+trHTML);
}
//条件检索预警信息
function  classifyWarnMessage(){
    var val_start =new Date($('#starttime').val());
    var val_end = new Date($('#endtime').val());
    var mystarttime = val_start.getTime()/1000;
    var myendtime = val_end.getTime()/1000;
    var mylevel = $('#levelid').val();
    getWarnList(mystarttime,myendtime,pageSize,pageNo,mylevel);
}

//预警模态框预警信息
function takeWarn(id,time,value,status,user,level,detail){
    $("#modal_id").val(id);
    $("#bjNum").html(id==undefined?'':id);
    $("#bjTime").html(time==undefined?'':time);
    $("#bjValue").html(value==undefined?'':value);
    if(status=='false'){
        $("#bjHadel").html("未处理");
        $("#bjHadel").css('color','red');
    }else{
        $("#bjHadel").html("已处理");
        $("#bjHadel").css('color','#48691b');
    }
    $("#bgPeople").html(user==undefined?'':user);
    $("#bgJb").html(level==undefined?'':level);
    $("#bgDetail").html(detail==undefined?'':detail);
    $("#myModal").modal('show');
}
//提交预警处理信息
$("#centerPoin").click(function () {
    var id= $("#modal_id").val();
    var contentTextare = $("#contentTextare").val();
    var data_m={
        alarm_id:id,
        content : contentTextare
    };
    obj.put('/api/v1/alarm/'+id,data_m,function (data) {
        alert('处理成功');
        $("#myModal").modal('hide');
        $("#contentTextare").val("");
        var level = $('#levelid').val();
        var val_start =new Date($('#starttime').val());
        var val_end = new Date($('#endtime').val());
        var mystarttime = val_start.getTime()/1000;
        var myendtime = val_end.getTime()/1000;
        getWarnList(mystarttime,myendtime,pageSize,pageNo2,level);
    },function (err) {
        console.log(err);
       alert('处理失败，请重新处理');
        $("#contentTextare").val("");
    })
});
//获取关注列表
function getFocus(pageSize,pageNo) {
    var data_m = {
        page_size:pageSize,
        page: pageNo
    };
    obj.get('/api/v1/favorite',data_m,function(data){
        var mypageNo3= data.current_page==null?0:data.current_page;
        var mypages3 = data.total_pages==null?0:data.total_pages;
        $("#pageMsg3").html(mypageNo3+"/"+mypages3);
        $("#curpage3").val(mypageNo3);
        $("#tolpage3").val(mypages3);
        $('#toPage3').val(mypageNo3);
        $("#talTher").html(data.total_records);
        var _data=data.data;
        $("#focusBody").html("");
        for(var i=0;i<_data.length;i++){
            var object=_data[i];
            showFocusList(object);
        }        
    },function(){
        return;
    })
}
//展示关注列表
function showFocusList(obj){
    var html = $("#focusBody").html();
    var trHTML = "<tr class=\"tbody-tr\">";
    trHTML += "<td class=\"tbody-td\">"+obj.id+"</td>";
    //trHTML += '<td class="tbody-td"><a href="cropCon.html?id='+obj.id+'" + "&orgId=" + obj.org_id >'+obj.corp_name+'</a></td>';
	
	trHTML += '<td class="tbody-td"><a class="tbody-td-a" href="cropCon.html?id='+obj.id+'&orgId='+obj.my_org_id+'">'+obj.corp_name+'</a></td>';
   
    trHTML += '<td class="tbody-td">'+obj.region+'</td>';
    trHTML += '</tr>';
    $("#focusBody").html(html+trHTML);
}
$("#refreshone").click(function(){
    getMessageList(jgid,pageSize,pageNo);
});
$("#refreshtwo").click(function(){
    var level = $('#levelid').val();
    var val_start =new Date($('#starttime').val());
    var val_end = new Date($('#endtime').val());
    var mystarttime = val_start.getTime()/1000;
    var myendtime = val_end.getTime()/1000;
    getWarnList(mystarttime,myendtime,pageSize,pageNo,level)
});
$("#refreshthre").click(function(){
    getFocus(pageSize,pageNo);
});

//点击树形结构
function zTreeOnClick(event, treeId, treeNode) {
    var treeObj = $.fn.zTree.getZTreeObj("treeDemo");
    treeObj.checkNode(treeNode, true, true);
    jgid = treeNode.id;
   // console.log("jgid = "+jgid+"name="+treeNode.name)
    var name = treeNode.name;
    var treeflag = treeNode.isParent;
    var la = treeNode.latitude;
    var lg = treeNode.longitude;
    gotoMap2(la,lg,jgid,treeflag,name);
    getMessageList(jgid,pageSize,pageNo);
	setMainCorpSelectId(jgid);
}
function zTreeOnAsyncSuccess(event, treeId, treeNode, msg) {
    alert(msg);
    setNode();
};

function getTree(){
    //获取文件树结构

    obj.get('/api/v1/orgs','', function(data) {
        $.fn.zTree.init($("#treeDemo"), setting, data.orgs);
        var treeObj = $.fn.zTree.getZTreeObj("treeDemo");
        $("#treeDemo_1_check").addClass('radio_true_full');
        
        treeObj.expandAll(true);
		SelectNode();

    });
}
//获取行政区域
function getBoundary(bd,a,b) {
	if(mapSource == "QQ") return;
    var bdary = new BMap.Boundary();
    bdary.get(bd, function(rs) {
        //            map.clearOverlays();        //清除地图覆盖物
        var count = rs.boundaries.length; //行政区域的点有多少个
        if (count === 0) {
            return;
        }
        var pointArray = [];
        for (var i = 0; i < count; i++) {
            var ply = new BMap.Polygon(rs.boundaries[i], {
                strokeWeight: 3,
                strokeColor: "red",
                fillOpacity: 0.00000000001
            }); //建立多边形覆盖物
            map.addOverlay(ply); //添加覆盖物
            ply.disableMassClear(); //设置覆盖物无法被清除
            pointArray = pointArray.concat(ply.getPath());
        }
    });
    if(levelflag==2){
        map.centerAndZoom(new BMap.Point(a,b), 11);
        map.enableScrollWheelZoom();
        map.panTo(new BMap.Point(a+0.001,b));
    }else if(levelflag==0 || levelflag == 1){
        map.centerAndZoom(new BMap.Point(a,b), 10);
        map.enableScrollWheelZoom();
        map.panTo(new BMap.Point(a+0.001,b));
    }
}

//获取机构下的尾矿库
function getAllMarkers(orgId) {
	if(findOldMap() == 1) return; // add by lwh at 2017-12-15
	// 清空原有的标注和标签信息
	markers = [];
	labels = [];
	if($("#maptype").val().indexOf("百度") >= 0)
	{
		mapSource = "BD";
		//鼠标移动监听
		map.addEventListener("mousemove", function(e) {
			$("#newL").html(e.point.lng + "," + e.point.lat);
		});
		getBoundary(boundary,lg,la);
	}
	else if($("#maptype").val().indexOf("腾讯") >= 0)
	{
		mapSource = "QQ";
		//鼠标移动监听
		qq.maps.event.addListener(map, "mousemove", function (e) {
		$("#newL").html(e.latLng.getLng().toFixed(6) + "," + e.latLng.getLat().toFixed(6));
		});
	}
    	
    var data = {
        org_id: orgId
    };
    obj.get('/api/v1/corps?page_size=1000', data, function(data) {
        var data = data.data;
		var cssC = {
			color: "black",
			fontSize: '12px',
			height: '20px',
			lineHeight: '20px',
			fontFamily: '微软雅黑',
			border: '1px solid #aaa', //原始为none，如此增加字体边框宽度和颜色
			backgroundColor: 'white'  //原始为none，切换字体底板颜色
		};
		try{
        for (var i = 0; i < data.length; i++) {
            var corp = data[i];
            
            var latitude = corp.latitude;
            var longitude = corp.longitude;
            var name = corp.corp_name;
			// 百度地图标注
			if(mapSource == "BD")
			{
				var point = new BMap.Point(longitude,latitude);
				var icon = new BMap.Icon('./images/icon/'+iconPng[corp.status%3], new BMap.Size(32, 32), {
						anchor: new BMap.Size(32, 32)});
				var label = new BMap.Label(name, {offset: new BMap.Size(20, -10)});
				label.setStyle(cssC);
				var marker = new BMap.Marker(point, { icon: icon });			
				marker.status = corp.status;
				marker.id = corp.id;
				marker.orgId = corp.my_org_id;
				markers.push(marker);
				labels.push(label);
			}
			// 腾讯地图标注
			else if(mapSource == "QQ")
			{
				var qqIcon = new qq.maps.MarkerImage(
				'./images/icon/'+iconPng[corp.status % 3],
				  new qq.maps.Size(32, 32), 
				  new qq.maps.Point(0, 0),
				  new qq.maps.Point(0, 0)
				);				
				var label = new qq.maps.Label({
                clickable: false, //如果为true，表示可点击，默认true                
                content: name, //标签的文本                
                map: map, //显示标签的地图
                offset: new qq.maps.Size(20, -10), //相对于position位置偏移值，x方向向右偏移为正值，y方向向下偏移为正值，反之为负
                position: new qq.maps.LatLng(latitude,longitude), //标签位置坐标，若offset不设置，默认标签左上角对准该位置                
                style: cssC, //Label样式                
                visible: true, //如果为true，表示标签可见，默认为true                
                //zIndex: 1000 //标签的z轴高度，zIndex大的标签，显示在zIndex小的前面
				});
				var marker = new qq.maps.Marker({position: new qq.maps.LatLng(latitude,longitude), map: map, icon: qqIcon});
				marker.status = corp.status;
				marker.id = corp.id;
				marker.orgId = corp.my_org_id;
				markers.push(marker);
				labels.push(label);
			}
        }
		if(mapSource == "BD")
		{
			if (!cluflag) {
				var lg= map.getCenter().lng;
				var la =map.getCenter().lat;
				var zoom = map.getZoom();
				closeClusterMarkers(lg,la,zoom);
				addMarker(markers,labels);
			} else if (cluflag) {
				clusterMarkers(markers);
				addMarker(markers,labels);
			}
		}
		else if(mapSource == "QQ")
		{
			if (!cluflag) {
				var lg= map.getCenter().lng;
				var la =map.getCenter().lat;
				var zoom = map.getZoom();
			} else if (cluflag) {
				clusterQmapMarkers(markers);				
			}
			addQmapMarker(markers,labels);
		}
	}catch(e){
		;//console.log(e);  返回时会出现QQ地图的map找不到对象，所以暂时强制浏览器重新加载页面
		location.reload(); 
	}
    }, function(err) {
        return;
    });

}

//显示聚合
function clusterMarkers(markers) {
    markerClusterer = new BMapLib.MarkerClusterer(map, {
        markers: markers,
        isAverangeCenter: true
    });
}

//获取每个聚合点的最高状态
// function getMarkerClusterersStatusMax(markerClusterer){
//     var mycluster = [];
//     var arr_status=[];//最高状态数组
//     var arr_clusters=markerClusterer._clusters;  //聚合点数组
//     for(var i=0;i<arr_clusters.length;i++){
//         var arr_markers=arr_clusters[i]._markers; //聚合点内的每个点
//         var arr_markersStatus=[];   //聚合点内的每个点的状态数组
//         for(var j=0;j<arr_markers.length;j++){
//             arr_markersStatus.push(arr_markers[j].status);
//         }
//         var max_status=arr_markersStatus.max();
//         arr_status.push(max_status);
//         mycluster[i]={
//             clustt : arr_clusters[i],
//             stt : max_status
//         };
//     }
//     // return arr_status;
//     return mycluster;
// }

//关闭聚合
function closeClusterMarkers(lg,la,zoom) {
    map = new BMap.Map("allmap");
    map.centerAndZoom(new BMap.Point(lg,la),parseInt(zoom));
    map.enableScrollWheelZoom();
    getBoundary(boundary);
    // //鼠标移动监听
    map.addEventListener("mousemove", function(e) {
        $("#newL").html(e.point.lng + "," + e.point.lat);
    });
}
//添加标注
function addMarker(markers, labels) {
    var i = 0;
    for (; i < markers.length; i++) {
        map.addOverlay(markers[i]);
        markers[i].setLabel(labels[i]);
        var mymarkadd = markers[i];
        mymarkadd.addEventListener("click",function(){
            location.href = "cropCon.html?id="+this.id + "&orgId=" + this.orgId;
        });
    }
    if(txtflag){
        showLabel(labels);
        myflage = 1;
    }else if(!txtflag){
        removeLabel(labels);
        myflage = 1;
    }
}

//去除标注中的label
function removeLabel(labels) {
    for (var i = 0; i < labels.length; i++) {
        labels[i].setStyle({
            display: 'none'
        })
    }
}
//显示标注中的label
function showLabel(labels) {
    for (var i = 0; i < labels.length; i++) {
        labels[i].setStyle({
            display: 'block'
        })
    }
}

// 放大地图
function getBiggerMap() {
    map.zoomTo(map.getZoom() + 1);
}
// 缩小地图
function getSmallMap() {
    map.zoomTo(map.getZoom() - 1);
}

//设置中心点
function setCenter() {
    if (confirm("设置当前位置为中心点")) {
        nowCenterLng = map.getCenter().lng;
        nowCenterLat = map.getCenter().lat;
    }
    else {
        return
    }
}
//返回中心点
function reterCenter() {
    nowCenter = new BMap.Point(nowCenterLng, nowCenterLat);
    map.panTo(nowCenter);
}

function cancelModal() {
    $("#myModal").modal('hide');
    $("#contentTextare").val("");
}

//复选参数
function changeParcon() {
    if ($("#wkItem").hasClass('parCon_all_defalut')) {
        $("#wkItem").removeClass('parCon_all_defalut');
        $("#wkItem").addClass('parCon_all_checked');
    } else {
        $("#wkItem").addClass('parCon_all_defalut');
        $("#wkItem").removeClass('parCon_all_checked');
    }
}

//测站名称切换
$("#nameSwitch").click(function() {
    if ($("#nameSwitch").hasClass('switchOn')) {
        $("#nameSwitch").removeClass('switchOn');
        $("#nameSwitch").addClass('switchOff');
        txtflag = 0;
        removeLabel(labels);
    } else {
        txtflag = 1;
        showLabel(labels);
        $("#nameSwitch").removeClass('switchOff');
        $("#nameSwitch").addClass('switchOn');
    }
});
//聚合开关切换
$("#clusterSwitch").click(function() {
    var lg= map.getCenter().lng;
    var la =map.getCenter().lat;
    var zoom = map.getZoom();
    if ($("#clusterSwitch").hasClass('switchOn')) {
        $("#clusterSwitch").removeClass('switchOn');
        $("#clusterSwitch").addClass('switchOff');
        cluflag = 0;
        if(mapSource == "BD")
		{
			closeClusterMarkers(lg,la,zoom);
			getAllMarkers();
			map.panTo(new BMap.Point(lg,la+0.001));
		}
		else if(mapSource == "QQ")
		{
			refreshQmapFlag = true;
			$('#maptype').trigger("change");
		}
    } else {
        cluflag = 1;
		if(mapSource == "BD")
		{
			clusterMarkers(markers);
			map.panTo(new BMap.Point(lg,la+0.001));
		}
		else if(mapSource == "QQ")
		{
			//map.setZoom(10);
			refreshQmapFlag = true;
			$('#maptype').trigger("change");
		}
        $("#clusterSwitch").removeClass('switchOff');
        $("#clusterSwitch").addClass('switchOn');
    }
});

//添加关注状态按钮
function addDiyDom(treeId, treeNode) {
    var aObj = $("#" + treeNode.tId + "_a");
    var editStr = '';
    if (treeNode.watched) {
        editStr = "<sapn class='focus_ico' id='focus_ico_" + treeNode.id + "' title='取消关注'></span>";
    } else {
        editStr = "<sapn class='unfocus_ico' id='focus_ico_" + treeNode.id + "' title='关注'></span>"
    }
    $(aObj).after(editStr);
    $("#focus_ico_" + treeNode.id).bind("click", function() {
        if ($(this).hasClass('focus_ico')) { //已经关注
            obj.delete('/api/v1/favorite/'+treeNode.id,function (data) {
                $("#focus_ico_" + treeNode.id).removeClass('focus_ico').addClass('unfocus_ico');
                $("#focus_ico_" + treeNode.id).attr('title', '关注');
                getTree();
                getFocus(pageSize,pageNo);
            },function (err) {
                console.log(err);
            });
        } else {
            obj.post('/api/v1/favorite/'+treeNode.id,"",function (data) {
                $("#focus_ico_" + treeNode.id).removeClass('unfocus_ico').addClass('focus_ico');
                $("#focus_ico_" + treeNode.id).attr('title', '取消关注');
                getTree();
                getFocus(pageSize,pageNo);
            },function (err) {
                console.log(err);
            });
        }
    })
}


/**
 * 翻页
 * @param page
 */
function goPage(page){
    var pageNo =parseInt($("#curpage").val());
    var pages = parseInt($("#tolpage").val());
    if (page == 'up') {
        if (pageNo == 1) {
            return;
        }
        pageNo = pageNo - 1;
    }
    if (page == 'down') {
        if (pageNo == pages || pageNo>pages) {
            return;
        }
        pageNo = pageNo + 1;
    }
    if (page == 'page') {
        if ($('#toPage').val() == '') {
            alert('请填入页数！');
            return;
        }
        if (isNaN($('#toPage').val())) {
            alert('请填入数字！');
            return;
        }
        if (parseInt($('#toPage').val()) > parseInt(pages)) {
            alert('填入页数过大！');
            return;
        }
        if (parseInt($('#toPage').val()) == pageNo) {
            return;
        } else {
            pageNo = $('#toPage').val();
        }
    }
    $("#pageNo").val(pageNo);
    getMessageList(jgid,pageSize,pageNo);
}

function goPage2(page){
    var pageNo2 =parseInt($("#curpage2").val());
    var pages2 = parseInt($("#tolpage2").val());
    if (page == 'up') {
        if (pageNo2 == 1) {
            return;
        }
        pageNo2 = pageNo2 - 1;
    }
    if (page == 'down') {
        if (pageNo2 == pages2 || pageNo2>pages2) {
            return;
        }
        pageNo2 = pageNo2 + 1;
    }
    if (page == 'page') {
        if ($('#toPage2').val() == '') {
            alert('请填入页数！');
            return;
        }
        if (isNaN($('#toPage2').val())) {
            alert('请填入数字！');
            return;
        }
        if (parseInt($('#toPage2').val()) > parseInt(pages2)) {
            alert('填入页数过大！');
            return;
        }
        if (parseInt($('#toPage2').val()) == pageNo2) {
            return;
        } else {
            pageNo2 = $('#toPage2').val();
        }
    }
    $("#pageNo2").val(pageNo2);
    var level = $('#levelid').val();
    var val_start =new Date($('#starttime').val());
    var val_end = new Date($('#endtime').val());
    var mystarttime = val_start.getTime()/1000;
    var myendtime = val_end.getTime()/1000;
    getWarnList(mystarttime,myendtime,pageSize,pageNo2,level);
}


function goPage3(page){
    var pageNo3 =parseInt($("#curpage3").val());
    var pages3 = parseInt($("#tolpage3").val());
    if (page == 'up') {
        if (pageNo3== 1) {
            return;
        }
        pageNo3 = pageNo3 - 1;
    }
    if (page == 'down') {
        if (pageNo3== pages3 || pageNo3>pages3) {
            return;
        }
        pageNo3 = pageNo3 + 1;
    }
    if (page == 'page') {
        if ($('#toPage3').val() == '') {
            alert('请填入页数！');
            return;
        }
        if (isNaN($('#toPage3').val())) {
            alert('请填入数字！');
            return;
        }
        if (parseInt($('#toPage3').val()) > parseInt(pages3)) {
            alert('填入页数过大！');
            return;
        }
        if (parseInt($('#toPage3').val()) == pageNo3) {
            return;
        } else {
            pageNo3 = $('#toPage3').val();
        }
    }
    $("#pageNo3").val(pageNo3);
    getFocus(pageSize,pageNo3);
}

//获取主页记录的地图选项 add by lwh at 2017-12-15
function getMainSelectMap(){
	var matchStr = "mainSelectMap";
	var arr,reg = new RegExp("(^| )"+matchStr+"=([^;]*)(;|$)");
	if(arr=document.cookie.match(reg))
		return unescape(arr[2]);
	else 
		return null;
}
//设置主页记录的地图选项 add by lwh at 2017-12-15
function setMainSelectMap(setVal){
	if(setVal == null || setVal == undefined) return;
	var exp = new Date();
	exp.setTime(exp.getTime()+360000);
	document.cookie = "mainSelectMap"+"="+setVal+";expires="+exp.toGMTString();
}
// add by lwh at 2017-11-25
//获取主页记录的被选企业ID
function getMainCorpSelectId(){
	var matchStr = "mainCorpSelectId";
	var arr,reg = new RegExp("(^| )"+matchStr+"=([^;]*)(;|$)");
	if(arr=document.cookie.match(reg))
		return unescape(arr[2]);
	else 
		return null;
}
//设置主页记录的被选企业ID
function setMainCorpSelectId(setVal){
	if(setVal == null || setVal == undefined) return;
	var exp = new Date();
	exp.setTime(exp.getTime()+360000);
	document.cookie = "mainCorpSelectId"+"="+setVal+";expires="+exp.toGMTString();
}
//清除主页记录的被选企业ID
function clearMainCorpSelectId(){
	var exp = new Date();
	exp.setTime(exp.getTime()-3600);
	document.cookie = "mainCorpSelectId"+"="+"00"+";expires="+exp.toGMTString();
}
// end by lwh

