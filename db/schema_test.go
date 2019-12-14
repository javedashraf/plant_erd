package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchema_ToErd(t *testing.T) {
	type fields struct {
		Tables []*Table
	}
	type args struct {
		showIndex bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "users and articles",
			fields: fields{
				Tables: []*Table{
					{
						Name: "articles",
						Columns: []*Column{
							{
								Name:       "id",
								Type:       "integer",
								NotNull:    true,
								PrimaryKey: true,
							},
							{
								Name:    "user_id",
								Type:    "integer",
								NotNull: true,
							},
						},
						ForeignKeys: []*ForeignKey{
							{
								FromColumn: "user_id",
								ToTable:    "users",
								ToColumn:   "id",
							},
						},
					},
					{
						Name: "users",
						Columns: []*Column{
							{
								Name:       "id",
								Type:       "integer",
								NotNull:    true,
								PrimaryKey: true,
							},
							{
								Name: "name",
								Type: "text",
							},
						},
					},
				},
			},
			args: args{
				showIndex: true,
			},
			want: `entity articles {
  * id : integer
  --
  * user_id : integer
}

entity users {
  * id : integer
  --
  name : text
}

articles }-- users`,
		},
		{
			name: "Reject foreign key which table isn't in schema",
			fields: fields{
				Tables: []*Table{
					{
						Name: "articles",
						Columns: []*Column{
							{
								Name:       "id",
								Type:       "integer",
								NotNull:    true,
								PrimaryKey: true,
							},
							{
								Name:    "user_id",
								Type:    "integer",
								NotNull: true,
							},
						},
						ForeignKeys: []*ForeignKey{
							{
								FromColumn: "user_id",
								ToTable:    "users",
								ToColumn:   "id",
							},
						},
					},
				},
			},
			args: args{
				showIndex: true,
			},
			want: `entity articles {
  * id : integer
  --
  * user_id : integer
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				Tables: tt.fields.Tables,
			}

			got := s.ToErd(tt.args.showIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSchema_SurroundingTablesWithin(t *testing.T) {
	tables := []*Table{
		{
			Name: "articles",
			Columns: []*Column{
				{
					Name:       "id",
					Type:       "integer",
					NotNull:    true,
					PrimaryKey: true,
				},
				{
					Name:    "user_id",
					Type:    "integer",
					NotNull: true,
				},
			},
			ForeignKeys: []*ForeignKey{
				{
					FromColumn: "user_id",
					ToTable:    "users",
					ToColumn:   "id",
				},
			},
		},
		{
			Name: "comments",
			Columns: []*Column{
				{
					Name:       "id",
					Type:       "integer",
					NotNull:    true,
					PrimaryKey: true,
				},
				{
					Name:    "article_id",
					Type:    "integer",
					NotNull: true,
				},
			},
			ForeignKeys: []*ForeignKey{
				{
					FromColumn: "article_id",
					ToTable:    "articles",
					ToColumn:   "id",
				},
			},
		},
		{
			Name: "followers",
			Columns: []*Column{
				{
					Name:       "id",
					Type:       "integer",
					NotNull:    true,
					PrimaryKey: true,
				},
				{
					Name:    "user_id",
					Type:    "integer",
					NotNull: true,
				},
				{
					Name:    "target_user_id",
					Type:    "integer",
					NotNull: true,
				},
			},
			ForeignKeys: []*ForeignKey{
				{
					FromColumn: "user_id",
					ToTable:    "users",
					ToColumn:   "id",
				},
				{
					FromColumn: "target_user_id",
					ToTable:    "users",
					ToColumn:   "id",
				},
			},
		},
		{
			Name: "followings",
			Columns: []*Column{
				{
					Name:       "id",
					Type:       "integer",
					NotNull:    true,
					PrimaryKey: true,
				},
				{
					Name:    "user_id",
					Type:    "integer",
					NotNull: true,
				},
				{
					Name:    "target_user_id",
					Type:    "integer",
					NotNull: true,
				},
			},
			ForeignKeys: []*ForeignKey{
				{
					FromColumn: "user_id",
					ToTable:    "users",
					ToColumn:   "id",
				},
				{
					FromColumn: "target_user_id",
					ToTable:    "users",
					ToColumn:   "id",
				},
			},
		},
		{
			Name: "likes",
			Columns: []*Column{
				{
					Name:    "article_id",
					Type:    "integer",
					NotNull: true,
				},
				{
					Name:    "user_id",
					Type:    "integer",
					NotNull: true,
				},
			},
			ForeignKeys: []*ForeignKey{
				{
					FromColumn: "article_id",
					ToTable:    "articles",
					ToColumn:   "id",
				},
				{
					FromColumn: "user_id",
					ToTable:    "users",
					ToColumn:   "id",
				},
			},
		},
		{
			Name: "revisions",
			Columns: []*Column{
				{
					Name:       "id",
					Type:       "integer",
					NotNull:    true,
					PrimaryKey: true,
				},
				{
					Name:    "article_id",
					Type:    "integer",
					NotNull: true,
				},
			},
			ForeignKeys: []*ForeignKey{
				{
					FromColumn: "article_id",
					ToTable:    "articles",
					ToColumn:   "id",
				},
			},
		},
		{
			Name: "users",
			Columns: []*Column{
				{
					Name:       "id",
					Type:       "integer",
					NotNull:    true,
					PrimaryKey: true,
				},
				{
					Name: "name",
					Type: "text",
				},
			},
		},
	}

	type fields struct {
		Tables []*Table
	}
	type args struct {
		tableName string
		distance  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "distance within 1 from articles",
			fields: fields{
				Tables: tables,
			},
			args: args{
				tableName: "articles",
				distance:  1,
			},
			want: []string{
				"articles",
				"comments",
				"likes",
				"revisions",
				"users",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				Tables: tt.fields.Tables,
			}

			got := s.SurroundingTablesWithin(tt.args.tableName, tt.args.distance)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSchema_Subset(t *testing.T) {
	articles := &Table{
		Name: "articles",
		Columns: []*Column{
			{
				Name:       "id",
				Type:       "integer",
				NotNull:    true,
				PrimaryKey: true,
			},
			{
				Name:    "user_id",
				Type:    "integer",
				NotNull: true,
			},
		},
		ForeignKeys: []*ForeignKey{
			{
				FromColumn: "user_id",
				ToTable:    "users",
				ToColumn:   "id",
			},
		},
	}

	comments := &Table{
		Name: "comments",
		Columns: []*Column{
			{
				Name:       "id",
				Type:       "integer",
				NotNull:    true,
				PrimaryKey: true,
			},
			{
				Name:    "article_id",
				Type:    "integer",
				NotNull: true,
			},
		},
		ForeignKeys: []*ForeignKey{
			{
				FromColumn: "article_id",
				ToTable:    "articles",
				ToColumn:   "id",
			},
		},
	}

	followers := &Table{
		Name: "followers",
		Columns: []*Column{
			{
				Name:       "id",
				Type:       "integer",
				NotNull:    true,
				PrimaryKey: true,
			},
			{
				Name:    "user_id",
				Type:    "integer",
				NotNull: true,
			},
			{
				Name:    "target_user_id",
				Type:    "integer",
				NotNull: true,
			},
		},
		ForeignKeys: []*ForeignKey{
			{
				FromColumn: "user_id",
				ToTable:    "users",
				ToColumn:   "id",
			},
			{
				FromColumn: "target_user_id",
				ToTable:    "users",
				ToColumn:   "id",
			},
		},
	}

	followings := &Table{
		Name: "followings",
		Columns: []*Column{
			{
				Name:       "id",
				Type:       "integer",
				NotNull:    true,
				PrimaryKey: true,
			},
			{
				Name:    "user_id",
				Type:    "integer",
				NotNull: true,
			},
			{
				Name:    "target_user_id",
				Type:    "integer",
				NotNull: true,
			},
		},
		ForeignKeys: []*ForeignKey{
			{
				FromColumn: "user_id",
				ToTable:    "users",
				ToColumn:   "id",
			},
			{
				FromColumn: "target_user_id",
				ToTable:    "users",
				ToColumn:   "id",
			},
		},
	}

	likes := &Table{
		Name: "likes",
		Columns: []*Column{
			{
				Name:    "article_id",
				Type:    "integer",
				NotNull: true,
			},
			{
				Name:    "user_id",
				Type:    "integer",
				NotNull: true,
			},
		},
		ForeignKeys: []*ForeignKey{
			{
				FromColumn: "article_id",
				ToTable:    "articles",
				ToColumn:   "id",
			},
			{
				FromColumn: "user_id",
				ToTable:    "users",
				ToColumn:   "id",
			},
		},
	}

	revisions := &Table{
		Name: "revisions",
		Columns: []*Column{
			{
				Name:       "id",
				Type:       "integer",
				NotNull:    true,
				PrimaryKey: true,
			},
			{
				Name:    "article_id",
				Type:    "integer",
				NotNull: true,
			},
		},
		ForeignKeys: []*ForeignKey{
			{
				FromColumn: "article_id",
				ToTable:    "articles",
				ToColumn:   "id",
			},
		},
	}

	users := &Table{
		Name: "users",
		Columns: []*Column{
			{
				Name:       "id",
				Type:       "integer",
				NotNull:    true,
				PrimaryKey: true,
			},
			{
				Name: "name",
				Type: "text",
			},
		},
	}

	tables := []*Table{
		articles,
		comments,
		followers,
		followings,
		likes,
		revisions,
		users,
	}

	type fields struct {
		Tables []*Table
	}
	type args struct {
		tableName string
		distance  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Schema
	}{
		{
			name: "distance within 1 from articles",
			fields: fields{
				Tables: tables,
			},
			args: args{
				tableName: "articles",
				distance:  1,
			},
			want: &Schema{
				Tables: []*Table{
					articles,
					comments,
					likes,
					revisions,
					users,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				Tables: tt.fields.Tables,
			}
			got := s.Subset(tt.args.tableName, tt.args.distance)
			assert.Equal(t, tt.want, got)
		})
	}
}
