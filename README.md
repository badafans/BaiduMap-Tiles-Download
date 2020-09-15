# BaiduMap-Tiles-Download
多线程百度地图瓦片图层下载工具

百度地图api接口 https://api.map.baidu.com/customimage/tile?&x=瓦片x&y=瓦片y&z=缩放层级&customid=地图风格

getid.exe 采用golang编写,用于百度BD09坐标系转墨卡托平面坐标系转百度瓦片ID

getpic.bat 负责单个进程下载百度瓦片图并校验是否下载成功

多线程百度地图瓦片下载.bat 负责调用getid.exe计算瓦片ID，提交下载任务给getpic.bat，并且负责管理curl进程(保护大量调用curl运行而不卡死系统)

仅供学习与参考,请勿用于商业用途
