package category

import (
	"github.com/go-funcards/category-service/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Category struct {
	CategoryID string    `json:"category_id" bson:"_id,omitempty"`
	OwnerID    string    `json:"owner_id" bson:"owner_id,omitempty"`
	BoardID    string    `json:"board_id" bson:"board_id,omitempty"`
	Name       string    `json:"name" bson:"name,omitempty"`
	Position   int32     `json:"position" bson:"position,omitempty"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at,omitempty"`
}

type Filter struct {
	CategoryIDs []string `json:"category_ids,omitempty"`
	OwnerIDs    []string `json:"owner_ids,omitempty"`
	BoardIDs    []string `json:"board_ids,omitempty"`
}

func (c Category) toProto() *v1.CategoriesResponse_Category {
	return &v1.CategoriesResponse_Category{
		CategoryId: c.CategoryID,
		OwnerId:    c.OwnerID,
		BoardId:    c.BoardID,
		Name:       c.Name,
		Position:   c.Position,
		CreatedAt:  timestamppb.New(c.CreatedAt),
	}
}

func CreateCategory(in *v1.CreateCategoryRequest) Category {
	return Category{
		CategoryID: in.GetCategoryId(),
		OwnerID:    in.GetOwnerId(),
		BoardID:    in.GetBoardId(),
		Name:       in.GetName(),
		Position:   in.GetPosition(),
		CreatedAt:  time.Now().UTC(),
	}
}

func UpdateCategory(in *v1.UpdateCategoryRequest) Category {
	return Category{
		CategoryID: in.GetCategoryId(),
		BoardID:    in.GetBoardId(),
		Name:       in.GetName(),
		Position:   in.GetPosition(),
	}
}

func CreateFilter(in *v1.CategoriesRequest) Filter {
	return Filter{
		CategoryIDs: in.GetCategoryIds(),
		OwnerIDs:    in.GetOwnerIds(),
		BoardIDs:    in.GetBoardIds(),
	}
}
