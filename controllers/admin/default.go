package admin

/**
 * 文件描述
 * Create on 11/28/17 2:17 PM
 * Create by hu
 */
import (
	"bloggo/models"
	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
	"strconv"
	//"fmt"
	"net/http"
	//"fmt"
	"os"
	"io"
	"log"
)


type SuccessController struct{
	BaseController
}


type IndexController struct {
	BaseController
}
type ArticleController struct {
	BaseController
}
type UploadController struct {
	BaseController
}


func (this *SuccessController) Get() {
	this.TplName="admin/success.html"
}


/**
* 首页
 */
func (this *IndexController) Get() {
	this.TplName = "admin/welcome.html"
	this.Layout = "admin/layout.html"

}


/**
*添加文章
 */
func (this *ArticleController) AddArticle() {
	this.TplName = "admin/article/add.html"
	this.Layout = "admin/layout.html"
	if this.Ctx.Request.Method == "POST" {
		art := new(models.Article)
		o := orm.NewOrm()
		if err := this.ParseForm(art); err != nil {
			seelog.Error("文章添加失败")
		}else{
			_, err := o.Insert(art)
			if err == nil {
				seelog.Info("文章添加成功")
				this.Success("提交成功","/admin/success",1)
			}else{
				this.Error("提交失败","",2)
			}
		}
	}
}

/**
文章列表
 */
func (this *ArticleController) ListArticle()  {
	o := orm.NewOrm()
	var info []models.Article
	nums,err :=o.QueryTable("article").All(&info)
	if err ==nil && nums >0 {
		this.Data["Total"] = len(info)
		this.Data["Pages"] = info

	}else{
		this.Data["Total"] = 0
		this.Data["Pages"] = 0
	}
	this.TplName = "admin/article/list.html"
	this.Layout = "admin/layout.html"
}


/**
文章详情
 */
func (this *ArticleController) GetDetailArticle()  {

	id := this.Ctx.Input.Param(":id")
	intid, err := strconv.Atoi(id)
	if err ==nil && intid >0{
		o := orm.NewOrm()
		var info []models.Article
		nums,err :=o.QueryTable("article").Filter("id",intid).All(&info)
		if err ==nil && nums >0 {
			//fmt.Println(info[0].Title)
			this.Data["Total"] = len(info)
			this.Data["Pages"] = info[0]

		}else{
			this.Data["Total"] = 0
			this.Data["Pages"] = 0
		}
	}
	this.TplName = "admin/article/detail.html"
	this.Layout = "admin/layout.html"
}

func (this *ArticleController) UpdateArticle()  {


	if this.Ctx.Request.Method == "GET" {
		id := this.Ctx.Input.Param(":id")
		intid, err := strconv.Atoi(id)
		if err ==nil && intid >0{
			o := orm.NewOrm()
			var info []models.Article
			nums,err :=o.QueryTable("article").Filter("id",intid).All(&info)
			if err ==nil && nums >0 {
				this.Data["Total"] = len(info)
				this.Data["Pages"] = info[0]

			}else{
				this.Data["Total"] = 0
				this.Data["Pages"] = 0
			}

			//this.TplName = "admin/article/update.html"
			//this.Layout = "admin/layout.html"
		}
	} else if this.Ctx.Request.Method == "POST" {
		articleid := this.GetString("articleid")
		title := this.GetString("title")
		content := this.GetString("content")
		if articleid != "" && title !="" && content !="" {
			o := orm.NewOrm()
			num, err := o.QueryTable("article").Filter("id", articleid).Update(orm.Params{
				"title": title,
				"content": content,
			})
			if err ==nil &&num >0{
						seelog.Info("文章添加成功")
						this.Success("提交成功","",1)
			}else{
				this.Error("提交失败","",2)
			}
		}

	}


	this.TplName = "admin/article/update.html"
	this.Layout = "admin/layout.html"

}


/**
*上传图片
 */
func (this *UploadController) PostUpload(w http.ResponseWriter, r *http.Request)  {
	//设置内存大小
	r.ParseMultipartForm(32 << 20);
	//获取上传的文件组
	files := r.MultipartForm.File["file"];
	len := len(files);
	for i := 0; i < len; i++ {
		//打开上传文件
		file, err := files[i].Open();
		defer file.Close();
		if err != nil {
			log.Fatal(err);
		}
		//创建上传目录
		os.Mkdir("./upload", os.ModePerm);
		//创建上传文件
		cur, err := os.Create("/static/upload/" + files[i].Filename);
		defer cur.Close();
		if err != nil {
			log.Fatal(err);
		}
		io.Copy(cur, file);
	}
}
