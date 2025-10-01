package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raihaninkam/finalPhase3/internals/models"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/raihaninkam/finalPhase3/internals/utils"
	"github.com/raihaninkam/finalPhase3/pkg"
)

type AuthHandler struct {
	ar *repositories.AuthRepositories
}

func NewAuthHandler(ar *repositories.AuthRepositories) *AuthHandler {
	return &AuthHandler{ar: ar}
}

// Register godoc
// @Summary     Register User
// @Description Daftar User baru dengan email dan password. Password akan di-hash sebelum disimpan.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       body body models.UserAuth true "Register Request"
// @Success     201 {object} map[string]interface{} "User berhasil didaftarkan"
// @Failure     400 {object} map[string]interface{} "Bad Request - Input tidak valid (email, password)"
// @Failure     409 {object} map[string]interface{} "Conflict - Email sudah terdaftar"
// @Failure     500 {object} map[string]interface{} "Internal Server Error"
// @Router      /auth/register [post]
func (a *AuthHandler) Register(ctx *gin.Context) {
	// menerima body
	var body models.AuthRequest
	if err := ctx.ShouldBind(&body); err != nil {
		if strings.Contains(err.Error(), "required") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Email dan Password harus diisi",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	// validasi email
	if err := utils.RegisterValidation(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// validasi password
	if err := utils.RegisterValidation(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// cek apakah email sudah terdaftar
	exists, err := a.ar.CheckEmailExists(ctx.Request.Context(), body.Email)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}
	if exists {
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Email sudah terdaftar",
		})
		return
	}

	// hash password sebelum disimpan
	hc := pkg.NewHashConfig()
	hc.UseRecommended() // menggunakan konfigurasi yang direkomendasikan
	hashedPassword, err := hc.GenHash(body.Password)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: hashedPassword,
	}

	// simpan user ke database
	err = a.ar.CreateAccount(ctx.Request.Context(), &user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "Email sudah terdaftar",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	// response sukses tanpa data user
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User berhasil didaftarkan",
	})
}

// Login godoc
// @Summary     Login User
// @Description Login dengan email dan password. Jika sukses, akan mengembalikan JWT Token untuk autentikasi.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       body body models.UserAuth true "Login Request"
// @Success     200 {object} map[string]interface{} "Berhasil login, kembalikan token"
// @Failure     400 {object} map[string]interface{} "Bad Request - Email/Password salah atau input tidak valid"
// @Failure     500 {object} map[string]interface{} "Internal Server Error"
// @Router      /auth/login [post]
func (a *AuthHandler) Login(ctx *gin.Context) {
	// menerima body
	var body models.User

	if err := ctx.ShouldBind(&body); err != nil {
		if strings.Contains(err.Error(), "required") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Email dan Password harus diisi",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	// validasi basic email format (opsional untuk login)
	if body.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Email harus diisi",
		})
		return
	}

	if body.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Password harus diisi",
		})
		return
	}

	// ambil data user dari database
	user, err := a.ar.GetEmail(ctx.Request.Context(), body.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Email atau password salah",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	// bandingkan password menggunakan hash function Anda
	hc := pkg.NewHashConfig()
	isMatched, err := hc.CompareHashAndPassword(body.Password, user.Password)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		// Cek jika error terkait dengan hash format atau crypto
		if strings.Contains(err.Error(), "hash") ||
			strings.Contains(err.Error(), "crypto") ||
			strings.Contains(err.Error(), "argon2id") ||
			strings.Contains(err.Error(), "format") {
			log.Println("Error during password hashing/comparison")
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	// jika password tidak cocok
	if !isMatched {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Email atau password salah",
		})
		return
	}

	// jika match, buat JWT token
	claims := pkg.NewJWTClaims(user.ID, "user")
	jwtToken, err := claims.GenToken()
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	// response sukses dengan token
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"token":   jwtToken,
	})
}

// UpdateProfile godoc
// @Summary     Update User Profile
// @Description Update profil user (name, bio, avatar). Avatar akan diupload jika disertakan.
// @Tags        Auth
// @Accept      multipart/form-data
// @Produce     json
// @Security    BearerAuth
// @Param       name formData string false "Nama user"
// @Param       bio formData string false "Bio user"
// @Param       avatar formData file false "Avatar image (JPG, PNG, max 2MB)"
// @Success     200 {object} map[string]interface{} "Profile berhasil diupdate"
// @Failure     400 {object} map[string]interface{} "Bad Request - Input tidak valid"
// @Failure     401 {object} map[string]interface{} "Unauthorized - Token tidak valid"
// @Failure     500 {object} map[string]interface{} "Internal Server Error"
// @Router      /auth/profile [put]
func (a *AuthHandler) UpdateProfile(ctx *gin.Context) {
	// Ambil userID dari context (dari middleware JWT)
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Ambil form field
	name := ctx.PostForm("name")
	bio := ctx.PostForm("bio")

	// Ambil file avatar jika ada
	file, err := ctx.FormFile("avatar")
	var avatarUrl string
	if err == nil {
		// Upload file avatar
		uploadedFile, err := utils.FileUpload(ctx, file, "avatar")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		avatarUrl = "/public/" + uploadedFile
	}

	// Validasi: minimal harus ada salah satu field yang diisi
	if name == "" && bio == "" && avatarUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Minimal satu field harus diisi (name, bio, atau avatar)",
		})
		return
	}

	// Buat object untuk update
	updateData := &models.UserUpdate{
		ID:        userID,
		Name:      name,
		Bio:       bio,
		AvatarUrl: avatarUrl,
	}

	// Update ke database
	updatedUser, err := a.ar.UpdateAccount(ctx.Request.Context(), updateData)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "User tidak ditemukan",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	// Response sukses
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile berhasil diupdate",
		"data": gin.H{
			"email":      updatedUser.Email,
			"name":       updatedUser.Name,
			"bio":        updatedUser.Bio,
			"avatar_url": updatedUser.AvatarUrl,
		},
	})
}
