package controller

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"
	"TinyTik/resp"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	resp.Response
	UserList []model.User `json:"user_list"`
}

type RelationActionResponse struct {
	resp.Response
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	redis := common.GetRedisClient()
	repo := repository.GetRelaRepo()
	if user, exist := redis.UserLoginInfo(token); exist {
		toUser, err := repository.NewUserRepository().GetUserById(to_user_id)
		if err != nil {
			c.JSON(http.StatusOK, RelationActionResponse{
				Response: resp.Response{
					// TODO 统一定义返回码
					StatusCode: -1,
					StatusMsg:  "请求关注用户失败，用户不存在",
				},
			})
			return
		}
		switch action_type {
		case 1: // TODO: 统一定义操作码，关注操作
			{
				// 如果没有关注的话，进行关注
				if !repo.Followed(&user, &toUser) {
					// FIXME 修改user和toUser时需要加锁
					user.FollowCount += 1
					toUser.FollowerCount += 1
					if err := repo.UpdateRelation(user, toUser, model.FOLLOW); err != nil {
						c.JSON(http.StatusOK, RelationActionResponse{
							Response: resp.Response{
								StatusCode: -1,
								StatusMsg:  "关注失败",
							},
						})
					} else {
						c.JSON(http.StatusOK, RelationActionResponse{
							Response: resp.Response{
								StatusCode: 0,
								StatusMsg:  "关注成功",
							},
						})
					}
				} else { // 已经关注，不能重复关注
					c.JSON(http.StatusOK, RelationActionResponse{
						Response: resp.Response{
							StatusCode: -1,
							StatusMsg:  "不能重复关注",
						},
					})
				}
			}
		case 2: // 取关操作
			{
				// 如果存在关注关系的话进行取关
				if repo.Followed(&user, &toUser) {
					user.FollowCount -= 1
					toUser.FollowerCount -= 1

					if err := repo.UpdateRelation(user, toUser, model.UNFOLLOW); err != nil {
						c.JSON(http.StatusOK, RelationActionResponse{
							Response: resp.Response{
								StatusCode: -1,
								StatusMsg:  "取关失败",
							},
						})
					} else {
						c.JSON(http.StatusOK, RelationActionResponse{
							Response: resp.Response{
								StatusCode: 0,
								StatusMsg:  "取关成功",
							},
						})
					}
				} else {
					c.JSON(http.StatusOK, RelationActionResponse{
						Response: resp.Response{
							StatusCode: -1,
							StatusMsg:  "没有关注，无法取关",
						},
					})
				}
			}
		default:
			{
				c.JSON(http.StatusOK, resp.Response{
					StatusCode: -1,
					StatusMsg:  "非法参数",
				})
			}
		}
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

type FollowListResponse struct {
	resp.Response
	Users []model.User `json:"user_list"`
}

type FollowerListResponse FollowListResponse
type FriendListResponse FollowListResponse

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	redis := common.GetRedisClient()
	if _, exist := redis.UserLoginInfo(token); exist {
		id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		repo := repository.GetRelaRepo()
		user, err := repository.NewUserRepository().GetUserById(id)
		if err != nil {
			fmt.Printf("用户%v不存在\n", id)
			c.JSON(http.StatusOK, FollowListResponse{
				Response: resp.Response{
					StatusCode: -1,
					StatusMsg:  "用户不存在",
				},
			})
			return
		}
		res, err := repo.GetFollowListById(user.Id)
		if err != nil {
			c.JSON(http.StatusOK, FollowListResponse{
				Response: resp.Response{
					StatusCode: -1,
					StatusMsg:  "获取关注列表失败",
				},
			})
			return
		}
		c.JSON(http.StatusOK, FollowListResponse{
			Response: resp.Response{
				StatusCode: http.StatusOK,
				StatusMsg:  "查询成功",
			},
			Users: res,
		})
	} else {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: resp.Response{
				StatusCode: -1,
				StatusMsg:  "Access Denied.",
			},
		})
	}
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	redis := common.GetRedisClient()
	if _, exist := redis.UserLoginInfo(token); exist {
		id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		repo := repository.GetRelaRepo()
		user, err := repository.NewUserRepository().GetUserById(id)
		if err != nil {
			fmt.Printf("用户%v不存在\n", id)
			c.JSON(http.StatusOK, FollowerListResponse{
				Response: resp.Response{
					StatusCode: -1,
					StatusMsg:  "用户不存在",
				},
			})
			return
		}
		res, err := repo.GetFollowerListById(user.Id)
		if err != nil {
			c.JSON(http.StatusOK, FollowerListResponse{
				Response: resp.Response{
					StatusCode: -1,
					StatusMsg:  "获取粉丝列表失败",
				},
			})
			return
		}
		c.JSON(http.StatusOK, FollowerListResponse{
			Response: resp.Response{
				StatusCode: http.StatusOK,
				StatusMsg:  "查询成功",
			},
			Users: res,
		})
	} else {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: resp.Response{
				StatusCode: -1,
				StatusMsg:  "Access Denied.",
			},
		})
	}
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	token := c.Query("token")
	redis := common.GetRedisClient()
	if _, exist := redis.UserLoginInfo(token); exist {
		id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		repo := repository.GetRelaRepo()
		user, err := repository.NewUserRepository().GetUserById(id)
		if err != nil {
			fmt.Printf("用户%v不存在\n", id)
			c.JSON(http.StatusOK, FriendListResponse{
				Response: resp.Response{
					StatusCode: -1,
					StatusMsg:  "用户不存在",
				},
			})
			return
		}
		res, err := repo.GetFriendListById(user.Id)
		if err != nil {
			c.JSON(http.StatusOK, FriendListResponse{
				Response: resp.Response{
					StatusCode: -1,
					StatusMsg:  "获取朋友列表失败",
				},
			})
			return
		}
		c.JSON(http.StatusOK, FriendListResponse{
			Response: resp.Response{
				StatusCode: http.StatusOK,
				StatusMsg:  "查询成功",
			},
			Users: res,
		})
	} else {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: resp.Response{
				StatusCode: -1,
				StatusMsg:  "Access Denied.",
			},
		})
	}
}
