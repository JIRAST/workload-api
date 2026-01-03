package handler

import (
	"context"
		"fmt"
			"net/http"
				"os"
					"time"

						"github.com/jackc/pgx/v5"
						)

						func Handler(w http.ResponseWriter, r *http.Request) {
							dbUrl := os.Getenv("DB_URL")
								ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
									defer cancel()

										// ตั้งค่า Header เพื่อรองรับการเรียกจากหน้าเว็บ (CORS)
											w.Header().Set("Content-Type", "application/json")
												w.Header().Set("Access-Control-Allow-Origin", "*")

													conn, err := pgx.Connect(ctx, dbUrl)
														if err != nil {
																w.WriteHeader(http.StatusInternalServerError)
																		fmt.Fprintf(w, `{"status": "error", "message": "Connection failed"}`)
																				return
																					}
																						defer conn.Close(ctx)

																							var version string
																								err = conn.QueryRow(ctx, "SELECT version()").Scan(&version)
																									if err != nil {
																											w.WriteHeader(http.StatusInternalServerError)
																													fmt.Fprintf(w, `{"status": "error", "message": "Query failed"}`)
																															return
																																}

																																	w.WriteHeader(http.StatusOK)
																																		fmt.Fprintf(w, `{"status": "ok", "database": "connected", "version": "%s"}`, version)
																																		}
																																		