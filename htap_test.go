package gosnowflake

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestSortingByPriority(t *testing.T) {
	qcc := (&queryContextCache{}).init()
	sc := htapTestSnowflakeConn()

	qceA := queryContextEntry{ID: 12, Timestamp: 123, Priority: 7, Context: "a"}
	qceB := queryContextEntry{ID: 13, Timestamp: 124, Priority: 9, Context: "b"}
	qceC := queryContextEntry{ID: 14, Timestamp: 125, Priority: 6, Context: "c"}
	qceD := queryContextEntry{ID: 15, Timestamp: 126, Priority: 8, Context: "d"}

	t.Run("Add to empty cache", func(t *testing.T) {
		qcc.add(sc, qceA)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA}) {
			t.Fatalf("no entries added to cache. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with different id, timestamp and priority - greater priority", func(t *testing.T) {
		qcc.add(sc, qceB)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with different id, timestamp and priority - lesser priority", func(t *testing.T) {
		qcc.add(sc, qceC)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceC, qceA, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with different id, timestamp and priority - priority in the middle", func(t *testing.T) {
		qcc.add(sc, qceD)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceC, qceA, qceD, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
}

func TestAddingQcesWithTheSameIdAndLaterTimestamp(t *testing.T) {
	qcc := (&queryContextCache{}).init()
	sc := htapTestSnowflakeConn()

	qceA := queryContextEntry{ID: 12, Timestamp: 123, Priority: 7, Context: "a"}
	qceB := queryContextEntry{ID: 13, Timestamp: 124, Priority: 9, Context: "b"}
	qceC := queryContextEntry{ID: 12, Timestamp: 125, Priority: 6, Context: "c"}
	qceD := queryContextEntry{ID: 12, Timestamp: 126, Priority: 6, Context: "d"}

	t.Run("Add to empty cache", func(t *testing.T) {
		qcc.add(sc, qceA)
		qcc.add(sc, qceB)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA, qceB}) {
			t.Fatalf("no entries added to cache. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with different priority", func(t *testing.T) {
		qcc.add(sc, qceC)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceC, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with same priority", func(t *testing.T) {
		qcc.add(sc, qceD)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceD, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
}

func TestAddingQcesWithTheSameIdAndSameTimestamp(t *testing.T) {
	qcc := (&queryContextCache{}).init()
	sc := htapTestSnowflakeConn()

	qceA := queryContextEntry{ID: 12, Timestamp: 123, Priority: 7, Context: "a"}
	qceB := queryContextEntry{ID: 13, Timestamp: 124, Priority: 9, Context: "b"}
	qceC := queryContextEntry{ID: 12, Timestamp: 123, Priority: 6, Context: "c"}
	qceD := queryContextEntry{ID: 12, Timestamp: 123, Priority: 6, Context: "d"}

	t.Run("Add to empty cache", func(t *testing.T) {
		qcc.add(sc, qceA)
		qcc.add(sc, qceB)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA, qceB}) {
			t.Fatalf("no entries added to cache. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with different priority", func(t *testing.T) {
		qcc.add(sc, qceC)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceC, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with same priority", func(t *testing.T) {
		qcc.add(sc, qceD)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceC, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
}

func TestAddingQcesWithTheSameIdAndEarlierTimestamp(t *testing.T) {
	qcc := (&queryContextCache{}).init()
	sc := htapTestSnowflakeConn()

	qceA := queryContextEntry{ID: 12, Timestamp: 123, Priority: 7, Context: "a"}
	qceB := queryContextEntry{ID: 13, Timestamp: 124, Priority: 9, Context: "b"}
	qceC := queryContextEntry{ID: 12, Timestamp: 122, Priority: 6, Context: "c"}
	qceD := queryContextEntry{ID: 12, Timestamp: 122, Priority: 7, Context: "d"}

	t.Run("Add to empty cache", func(t *testing.T) {
		qcc.add(sc, qceA)
		qcc.add(sc, qceB)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with different priority", func(t *testing.T) {
		qcc.add(sc, qceC)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with same priority", func(t *testing.T) {
		qcc.add(sc, qceD)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
}

func TestAddingQcesWithDifferentId(t *testing.T) {
	qcc := (&queryContextCache{}).init()
	sc := htapTestSnowflakeConn()

	qceA := queryContextEntry{ID: 12, Timestamp: 123, Priority: 7, Context: "a"}
	qceB := queryContextEntry{ID: 13, Timestamp: 124, Priority: 9, Context: "b"}
	qceC := queryContextEntry{ID: 14, Timestamp: 122, Priority: 7, Context: "c"}
	qceD := queryContextEntry{ID: 15, Timestamp: 122, Priority: 6, Context: "d"}

	t.Run("Add to empty cache", func(t *testing.T) {
		qcc.add(sc, qceA)
		qcc.add(sc, qceB)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceA, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with same priority", func(t *testing.T) {
		qcc.add(sc, qceC)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceC, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
	t.Run("Add another entry with different priority", func(t *testing.T) {
		qcc.add(sc, qceD)
		if !reflect.DeepEqual(qcc.entries, []queryContextEntry{qceD, qceC, qceB}) {
			t.Fatalf("unexpected qcc entries. %v", qcc.entries)
		}
	})
}

func TestAddingQueryContextCacheEntry(t *testing.T) {
	runSnowflakeConnTest(t, func(sct *SCTest) {
		t.Run("First query (may be on empty cache)", func(t *testing.T) {
			entriesBefore := make([]queryContextEntry, len(sct.sc.queryContextCache.entries))
			copy(entriesBefore, sct.sc.queryContextCache.entries)
			sct.mustQuery("SELECT 1", nil)
			entriesAfter := sct.sc.queryContextCache.entries

			if !containsNewEntries(entriesAfter, entriesBefore) {
				t.Error("no new entries added to the query context cache")
			}
		})

		t.Run("Second query (cache should not be empty)", func(t *testing.T) {
			entriesBefore := make([]queryContextEntry, len(sct.sc.queryContextCache.entries))
			copy(entriesBefore, sct.sc.queryContextCache.entries)
			if len(entriesBefore) == 0 {
				t.Fatalf("cache should not be empty after first query")
			}
			sct.mustQuery("SELECT 2", nil)
			entriesAfter := sct.sc.queryContextCache.entries

			if !containsNewEntries(entriesAfter, entriesBefore) {
				t.Error("no new entries added to the query context cache")
			}
		})
	})
}

func containsNewEntries(entriesAfter []queryContextEntry, entriesBefore []queryContextEntry) bool {
	if len(entriesAfter) > len(entriesBefore) {
		return true
	}

	for _, entryAfter := range entriesAfter {
		for _, entryBefore := range entriesBefore {
			if !reflect.DeepEqual(entryBefore, entryAfter) {
				return true
			}
		}
	}

	return false
}

func TestPruneBySessionValue(t *testing.T) {
	qce1 := queryContextEntry{1, 1, 1, ""}
	qce2 := queryContextEntry{2, 2, 2, ""}
	qce3 := queryContextEntry{3, 3, 3, ""}

	testcases := []struct {
		size     string
		expected []queryContextEntry
	}{
		{
			size:     "1",
			expected: []queryContextEntry{qce1},
		},
		{
			size:     "2",
			expected: []queryContextEntry{qce1, qce2},
		},
		{
			size:     "3",
			expected: []queryContextEntry{qce1, qce2, qce3},
		},
		{
			size:     "4",
			expected: []queryContextEntry{qce1, qce2, qce3},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.size, func(t *testing.T) {
			sc := &snowflakeConn{
				cfg: &Config{
					Params: map[string]*string{
						queryContextCacheSizeParamName: &tc.size,
					},
				},
			}

			qcc := (&queryContextCache{}).init()

			qcc.add(sc, qce1)
			qcc.add(sc, qce2)
			qcc.add(sc, qce3)

			if !reflect.DeepEqual(qcc.entries, tc.expected) {
				t.Errorf("unexpected cache entries. expected: %v, got: %v", tc.expected, qcc.entries)
			}
		})
	}
}

func TestPruneByDefaultValue(t *testing.T) {
	qce1 := queryContextEntry{1, 1, 1, ""}
	qce2 := queryContextEntry{2, 2, 2, ""}
	qce3 := queryContextEntry{3, 3, 3, ""}
	qce4 := queryContextEntry{4, 4, 4, ""}
	qce5 := queryContextEntry{5, 5, 5, ""}
	qce6 := queryContextEntry{6, 6, 6, ""}

	sc := &snowflakeConn{
		cfg: &Config{
			Params: map[string]*string{},
		},
	}

	qcc := (&queryContextCache{}).init()
	qcc.add(sc, qce1)
	qcc.add(sc, qce2)
	qcc.add(sc, qce3)
	qcc.add(sc, qce4)
	qcc.add(sc, qce5)

	if len(qcc.entries) != 5 {
		t.Fatalf("Expected 5 elements, got: %v", len(qcc.entries))
	}

	qcc.add(sc, qce6)
	if len(qcc.entries) != 5 {
		t.Fatalf("Expected 5 elements, got: %v", len(qcc.entries))
	}
}

func TestNoQcesClearsCache(t *testing.T) {
	qce1 := queryContextEntry{1, 1, 1, ""}

	sc := &snowflakeConn{
		cfg: &Config{
			Params: map[string]*string{},
		},
	}

	qcc := (&queryContextCache{}).init()
	qcc.add(sc, qce1)

	if len(qcc.entries) != 1 {
		t.Fatalf("improperly inited cache")
	}

	qcc.add(sc)

	if len(qcc.entries) != 0 {
		t.Errorf("after adding empty context list cache should be cleared")
	}
}

func htapTestSnowflakeConn() *snowflakeConn {
	return &snowflakeConn{
		cfg: &Config{
			Params: map[string]*string{},
		},
	}
}

func TestQueryContextCacheDisabled(t *testing.T) {
	origDsn := dsn
	defer func() {
		dsn = origDsn
	}()
	dsn += "&disableQueryContextCache=true"
	runSnowflakeConnTest(t, func(sct *SCTest) {
		sct.mustExec("SELECT 1", nil)
		if len(sct.sc.queryContextCache.entries) > 0 {
			t.Error("should not contain any entries")
		}
	})
}

func TestHybridTablesE2E(t *testing.T) {
	if runningOnGithubAction() && !runningOnAWS() {
		t.Skip("HTAP is enabled only on AWS")
	}
	runID := time.Now().UnixMilli()
	testDb1 := fmt.Sprintf("hybrid_db_test_%v", runID)
	testDb2 := fmt.Sprintf("hybrid_db_test_%v_2", runID)
	runSnowflakeConnTest(t, func(sct *SCTest) {
		dbQuery := sct.mustQuery("SELECT CURRENT_DATABASE()", nil)
		defer dbQuery.Close()
		currentDb := make([]driver.Value, 1)
		dbQuery.Next(currentDb)
		defer func() {
			sct.mustExec(fmt.Sprintf("USE DATABASE %v", currentDb[0]), nil)
			sct.mustExec(fmt.Sprintf("DROP DATABASE IF EXISTS %v", testDb1), nil)
			sct.mustExec(fmt.Sprintf("DROP DATABASE IF EXISTS %v", testDb2), nil)
		}()

		t.Run("Run tests on first database", func(t *testing.T) {
			sct.mustExec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", testDb1), nil)
			sct.mustExec("CREATE HYBRID TABLE test_hybrid_table (id INT PRIMARY KEY, text VARCHAR)", nil)

			sct.mustExec("INSERT INTO test_hybrid_table VALUES (1, 'a')", nil)
			rows := sct.mustQuery("SELECT * FROM test_hybrid_table", nil)
			defer rows.Close()
			row := make([]driver.Value, 2)
			rows.Next(row)
			if row[0] != "1" || row[1] != "a" {
				t.Errorf("expected 1, got %v and expected a, got %v", row[0], row[1])
			}

			sct.mustExec("INSERT INTO test_hybrid_table VALUES (2, 'b')", nil)
			rows2 := sct.mustQuery("SELECT * FROM test_hybrid_table", nil)
			defer rows2.Close()
			rows2.Next(row)
			if row[0] != "1" || row[1] != "a" {
				t.Errorf("expected 1, got %v and expected a, got %v", row[0], row[1])
			}
			rows2.Next(row)
			if row[0] != "2" || row[1] != "b" {
				t.Errorf("expected 2, got %v and expected b, got %v", row[0], row[1])
			}
			if len(sct.sc.queryContextCache.entries) != 2 {
				t.Errorf("expected two entries in query context cache, got: %v", sct.sc.queryContextCache.entries)
			}
		})
		t.Run("Run tests on second database", func(t *testing.T) {
			sct.mustExec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", testDb2), nil)
			sct.mustExec("CREATE HYBRID TABLE test_hybrid_table_2 (id INT PRIMARY KEY, text VARCHAR)", nil)
			sct.mustExec("INSERT INTO test_hybrid_table_2 VALUES (3, 'c')", nil)

			rows := sct.mustQuery("SELECT * FROM test_hybrid_table_2", nil)
			defer rows.Close()
			row := make([]driver.Value, 2)
			rows.Next(row)
			if row[0] != "3" || row[1] != "c" {
				t.Errorf("expected 3, got %v and expected c, got %v", row[0], row[1])
			}
			if len(sct.sc.queryContextCache.entries) != 3 {
				t.Errorf("expected three entries in query context cache, got: %v", sct.sc.queryContextCache.entries)
			}
		})
		t.Run("Run tests on first database again", func(t *testing.T) {
			sct.mustExec(fmt.Sprintf("USE DATABASE %v", testDb1), nil)

			sct.mustExec("INSERT INTO test_hybrid_table VALUES (4, 'd')", nil)

			rows := sct.mustQuery("SELECT * FROM test_hybrid_table", nil)
			defer rows.Close()
			if len(sct.sc.queryContextCache.entries) != 3 {
				t.Errorf("expected three entries in query context cache, got: %v", sct.sc.queryContextCache.entries)
			}
		})
	})
}